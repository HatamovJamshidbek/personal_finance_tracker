package postgres

import (
	email2 "auth_service/api/email"
	"auth_service/api/models"
	pb "auth_service/genproto/auth_service"
	"auth_service/pkg/helper"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewUserRepo(db *pgxpool.Pool, log logger.ILogger) storage.IUserStorage {
	return &UserRepo{
		db:  db,
		log: log,
	}
}

func (repo *UserRepo) RegisterUser(ctx context.Context, request models.Register) (*models.RegisterResponse, error) {
	var (
		err      error
		query    string
		response models.RegisterResponse
		id       = uuid.New()
		timeNow  = time.Now()
	)

	if repo.db == nil {
		return nil, fmt.Errorf("this error is connection is db")
	}
	query = `INSERT INTO users(
                 id,
                 full_name,
                 username,
                 email,
                 phone,
                 password_hash,
                 image,
                 role,
                 created_at
             ) VALUES 
                 ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING 
                 id,
                 full_name,
                 username,
                 email,
                 phone,
                 password_hash,
                 image,
                 role,
                 created_at`
	err = repo.db.QueryRow(ctx, query,
		id,
		request.FullName,
		request.Username,
		request.Email,
		request.Phone,
		request.PasswordHash,
		request.Image,
		request.Role,
		timeNow).Scan(
		&response.Id,
		&response.FullName,
		&response.UserName,
		&response.Email,
		&response.Phone,
		&response.PasswordHash,
		&response.Image,
		&response.Role,
		&response.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repo *UserRepo) Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error) {
	var (
		err      error
		query    string
		response models.LoginResponse
	)
	fmt.Println("email", request.Email)
	fmt.Println("password", request.PasswordHash)

	query = `SELECT id,
       username,
       email,
       password_hash,
       role 
       FROM  users where
                       email=$1 and password_hash=$2
                       and deleted_at is null`
	err = repo.db.QueryRow(ctx, query, request.Email, request.PasswordHash).
		Scan(
			&response.Id,
			&response.UserName,
			&response.Email,
			&response.PasswordHash,
			&response.Role)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
func (repo *UserRepo) ResetPassword(ctx context.Context, request *models.LoginRequest) (string, error) {
	var (
		query string
		err   error
		id    string
	)
	_, err = email2.Email(request.Email)
	if err != nil {
		return "", err
	}

	query = `UPDATE users SET password_hash=$1 WHERE email=$2 and deleted_at is null RETURNING id`
	err = repo.db.QueryRow(ctx, query, request.Email, request.PasswordHash, request.Email).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (repo *UserRepo) UpdateUserProfile(ctx context.Context, request *pb.User) (*pb.User, error) {
	var (
		err      error
		response pb.User
		params   = make(map[string]interface{})
		query    = "update users set "
		filter   = ""
	)
	params["id"] = request.GetId()

	if request.Email != "" {

		params["email"] = request.GetEmail()

		filter += " email = @email, "

	}
	if request.Username != "" {

		params["user_name"] = request.GetUsername()

		filter += " user_name = @user_name, "

	}

	if request.Image != "" {

		params["image"] = request.GetImage()

		filter += " image = @image, "

	}
	if request.FullName != "" {

		params["full_name"] = request.GetFullName()

		filter += " full_name = @full_name, "

	}
	if request.Role != "" {

		params["role"] = request.GetRole()

		filter += " role = @role, "

	}
	if request.Phone != "" {

		params["phone"] = request.GetPhone()

		filter += " phone = @phone, "

	}
	query += filter + ` updated_at = now() where id = @id returning id, username,full_name,email, image,role,phone, created_at::text, updated_at::text `

	fullQuery, args := helper.ReplaceQueryParams(query, params)

	err = repo.db.QueryRow(ctx, fullQuery, args...).Scan(&response.Id, &response.Username, &response.FullName, &response.Email, &response.Image, &response.Role, &response.Phone, &response.CreatedAt, &response.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repo *UserRepo) GetUserProfile(ctx context.Context, request *pb.PrimaryKeyUser) (*pb.User, error) {
	var (
		err      error
		query    string
		response pb.User
	)

	fmt.Printf("Fetching profile for user ID: %v\n", request.GetId())

	query = `SELECT id, username, full_name, email, image, role, phone, created_at, updated_at 
             FROM users WHERE id=$1`

	// Use time.Time for timestamp columns
	var createdAt, updatedAt time.Time

	err = repo.db.QueryRow(ctx, query, request.GetId()).Scan(
		&response.Id, &response.Username, &response.FullName, &response.Email,
		&response.Image, &response.Role, &response.Phone, &createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %v not found", request.GetId())
		}
		return nil, fmt.Errorf("error fetching user profile: %v", err)
	}
	fmt.Println("User created at", createdAt)
	// Convert time.Time to Protobuf Timestamp
	response.CreatedAt = createdAt.String()
	response.UpdatedAt = updatedAt.String()

	return &response, nil
}
