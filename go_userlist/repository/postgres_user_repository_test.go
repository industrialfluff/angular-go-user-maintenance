package repository_test

import (
	"database/sql"
	"errors"
	"go_userlist/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostgresUserRepository", func() {
	var (
		repo *repository.PostgresUserRepository
		mock sqlmock.Sqlmock
		db   *sql.DB
	)

	BeforeEach(func() {
		var err error
		// Create a mock database connection
		db, mock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		// Inject the mock db into the repository
		repo, err = repository.NewPostgresUserRepository(db)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		repo.Close()
	})

	// Test cases
	Context("CreateUser", func() {
		It("should create a user successfully", func() {
			user := &repository.User{
				User_name:   "johndoe",
				First_name:  "John",
				Last_name:   "Doe",
				Email:       "john.doe@example.com",
				User_status: "active",
				Department:  "IT",
			}

			mock.ExpectQuery(`INSERT INTO public\.users`).
				WithArgs(user.User_name, user.First_name, user.Last_name, user.Email, user.User_status, user.Department).
				WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

			err := repo.CreateUser(user)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.User_id).To(Equal(1))
		})

		It("should return an error if insertion fails", func() {
			user := &repository.User{
				User_name:   "johndoe",
				First_name:  "John",
				Last_name:   "Doe",
				Email:       "john.doe@example.com",
				User_status: "active",
				Department:  "IT",
			}

			mock.ExpectQuery(`INSERT INTO public\.users`).
				WithArgs(user.User_name, user.First_name, user.Last_name, user.Email, user.User_status, user.Department).
				WillReturnError(errors.New("insert error"))

			err := repo.CreateUser(user)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetUserByID", func() {
		It("should return the correct user", func() {
			userID := 1
			mock.ExpectQuery(`SELECT (.+) FROM public\.users`).
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
					AddRow(1, "johndoe", "John", "Doe", "john.doe@example.com", "active", "IT"))

			user, err := repo.GetUserByID(userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			Expect(user.User_id).To(Equal(1))
			Expect(user.User_name).To(Equal("johndoe"))
		})

		It("should return ErrUserNotFound if no user is found", func() {
			userID := 999
			mock.ExpectQuery(`SELECT (.+) FROM public\.users`).
				WithArgs(userID).
				WillReturnError(sql.ErrNoRows)

			user, err := repo.GetUserByID(userID)
			Expect(err).To(Equal(repository.ErrUserNotFound))
			Expect(user).To(BeNil())
		})
	})

	Context("UpdateUser", func() {
		It("should update the user successfully", func() {
			user := &repository.User{
				User_id:     1,
				User_name:   "johndoe",
				First_name:  "John",
				Last_name:   "Doe",
				Email:       "john.doe@example.com",
				User_status: "active",
				Department:  "IT",
			}

			mock.ExpectExec(`UPDATE public\.users`).
				WithArgs(user.User_name, user.First_name, user.Last_name, user.Email, user.User_status, user.Department, user.User_id).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.UpdateUser(user)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if update fails", func() {
			user := &repository.User{
				User_id:     1,
				User_name:   "johndoe",
				First_name:  "John",
				Last_name:   "Doe",
				Email:       "john.doe@example.com",
				User_status: "active",
				Department:  "IT",
			}

			mock.ExpectExec(`UPDATE public\.users`).
				WithArgs(user.User_name, user.First_name, user.Last_name, user.Email, user.User_status, user.Department, user.User_id).
				WillReturnError(errors.New("update error"))

			err := repo.UpdateUser(user)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("DeleteUserByID", func() {
		It("should delete the user successfully", func() {
			userID := 1

			mock.ExpectQuery(`SELECT (.+) FROM public\.users`).
				WithArgs(userID).
				WillReturnRows(sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
					AddRow(1, "johndoe", "John", "Doe", "john.doe@example.com", "active", "IT"))

			mock.ExpectExec(`DELETE FROM public\.users`).
				WithArgs(userID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.DeleteUserByID(userID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return ErrUserNotFound if user doesn't exist", func() {
			userID := 999

			mock.ExpectQuery(`SELECT (.+) FROM public\.users`).
				WithArgs(userID).
				WillReturnError(sql.ErrNoRows)

			err := repo.DeleteUserByID(userID)
			Expect(err).To(Equal(repository.ErrUserNotFound))
		})
	})

})

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PostgresUserRepository Suite")
}
