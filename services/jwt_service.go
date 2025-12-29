package services

import (
	"errors"
	"fmt"
	"jwt-viewer/models"
	"jwt-viewer/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

// DecodeToken decodes a JWT without verifying the signature
func (s *JWTService) DecodeToken(tokenString string) (*models.DecodeResponse, error) {
	// Split the token into its three parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return &models.DecodeResponse{
			Error: "Invalid JWT format. Expected 3 parts separated by dots.",
		}, errors.New("invalid JWT format")
	}

	// Decode header
	headerBytes, err := utils.DecodeBase64URL(parts[0])
	if err != nil {
		return &models.DecodeResponse{
			Error: fmt.Sprintf("Failed to decode header: %v", err),
		}, err
	}

	header, err := utils.ParseJSONToClaims(headerBytes)
	if err != nil {
		return &models.DecodeResponse{
			Error: fmt.Sprintf("Failed to parse header JSON: %v", err),
		}, err
	}

	// Decode payload
	payloadBytes, err := utils.DecodeBase64URL(parts[1])
	if err != nil {
		return &models.DecodeResponse{
			Error: fmt.Sprintf("Failed to decode payload: %v", err),
		}, err
	}

	payload, err := utils.ParseJSONToClaims(payloadBytes)
	if err != nil {
		return &models.DecodeResponse{
			Error: fmt.Sprintf("Failed to parse payload JSON: %v", err),
		}, err
	}

	// Signature is already in base64url format
	signature := parts[2]

	return &models.DecodeResponse{
		Header:    header,
		Payload:   payload,
		Signature: signature,
	}, nil
}

// EncodeToken generates a new JWT with the given header, payload, and secret
func (s *JWTService) EncodeToken(req *models.EncodeRequest) (*models.EncodeResponse, error) {
	// Default to HS256 if no algorithm specified
	algorithm := req.Algorithm
	if algorithm == "" {
		algorithm = "HS256"
	}

	// Determine signing method
	var signingMethod jwt.SigningMethod
	switch algorithm {
	case "HS256":
		signingMethod = jwt.SigningMethodHS256
	case "HS384":
		signingMethod = jwt.SigningMethodHS384
	case "HS512":
		signingMethod = jwt.SigningMethodHS512
	default:
		return &models.EncodeResponse{
			Error: fmt.Sprintf("Unsupported algorithm: %s. Supported: HS256, HS384, HS512", algorithm),
		}, errors.New("unsupported algorithm")
	}

	// Create token
	token := jwt.NewWithClaims(signingMethod, jwt.MapClaims(req.Payload))

	// Override header if custom values are provided
	if req.Header != nil {
		for k, v := range req.Header {
			token.Header[k] = v
		}
	}

	// Ensure the algorithm in header matches
	token.Header["alg"] = algorithm

	// Sign the token
	tokenString, err := token.SignedString([]byte(req.Secret))
	if err != nil {
		return &models.EncodeResponse{
			Error: fmt.Sprintf("Failed to sign token: %v", err),
		}, err
	}

	return &models.EncodeResponse{
		Token: tokenString,
	}, nil
}

// VerifyToken verifies a JWT signature and checks its validity
func (s *JWTService) VerifyToken(req *models.VerifyRequest) (*models.VerifyResponse, error) {
	// Parse and verify the token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(req.Secret), nil
	})

	if err != nil {
		// Check specific error types
		if errors.Is(err, jwt.ErrTokenExpired) {
			claims, _ := token.Claims.(jwt.MapClaims)
			return &models.VerifyResponse{
				Valid:   false,
				Message: "Token is expired",
				Claims:  claims,
			}, nil
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return &models.VerifyResponse{
				Valid:   false,
				Message: "Token is not valid yet (nbf claim)",
			}, nil
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return &models.VerifyResponse{
				Valid:   false,
				Message: "Token is malformed",
				Error:   err.Error(),
			}, nil
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return &models.VerifyResponse{
				Valid:   false,
				Message: "Signature verification failed",
			}, nil
		}

		return &models.VerifyResponse{
			Valid:   false,
			Message: "Token validation failed",
			Error:   err.Error(),
		}, nil
	}

	// Token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &models.VerifyResponse{
			Valid:   false,
			Message: "Failed to parse claims",
		}, errors.New("failed to parse claims")
	}

	// Check if token is expired (additional check)
	isExpired := utils.IsTokenExpired(claims["exp"])
	message := "Token is valid"
	if isExpired {
		message = "Token signature is valid but token is expired"
	}

	return &models.VerifyResponse{
		Valid:   token.Valid,
		Message: message,
		Claims:  claims,
	}, nil
}

// ExtractClaimInfo extracts standard JWT claims and formats them
func (s *JWTService) ExtractClaimInfo(claims map[string]interface{}) *models.ClaimInfo {
	info := &models.ClaimInfo{
		Iss: utils.GetStringClaim(claims, "iss"),
		Aud: utils.GetStringClaim(claims, "aud"),
		Sub: utils.GetStringClaim(claims, "sub"),
	}

	// Parse time-based claims
	if exp := utils.ParseTimestamp(claims["exp"]); exp != nil {
		info.Exp = exp
		info.ExpString = utils.FormatTime(exp)
		info.IsExpired = utils.IsTokenExpired(claims["exp"])
	}

	if iat := utils.ParseTimestamp(claims["iat"]); iat != nil {
		info.Iat = iat
		info.IatString = utils.FormatTime(iat)
	}

	if nbf := utils.ParseTimestamp(claims["nbf"]); nbf != nil {
		info.Nbf = nbf
		info.NbfString = utils.FormatTime(nbf)
	}

	return info
}
