package repository

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestGetLoginData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock *sql.DB
	mockDB := NewMockDB(ctrl)

	// Create a Repository instance with the mockDB
	repo := NewRepository(mockDB)

	// Set up test input
	ctx := context.Background()
	input := GetLoginDataInput{PhoneNumber: "123456789"}

	// Set up expectations for the QueryRowContext and Scan methods
	expectedUserID := 1
	expectedFullName := "John Doe"
	expectedHashedPassword := "hashedpassword"

	mockDB.EXPECT().QueryRowContext(ctx, "SELECT id, full_name, password FROM users WHERE phone_number = $1", input.PhoneNumber).
		Return(mocks.NewMockRow(ctrl)).
		Times(1)

	mockRow := mocks.NewMockRow(ctrl)
	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
		// Simulate scanning data into the output parameters
		dest[0] = expectedUserID
		dest[1] = expectedFullName
		dest[2] = expectedHashedPassword
		return nil
	})

	// Call the GetLoginData function
	output, err := repo.GetLoginData(ctx, input)

	// Verify the result and error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output.UserID != expectedUserID {
		t.Errorf("Expected UserID %d, got %d", expectedUserID, output.UserID)
	}

	if output.FullName != expectedFullName {
		t.Errorf("Expected FullName %s, got %s", expectedFullName, output.FullName)
	}

	if output.HashedPassword != expectedHashedPassword {
		t.Errorf("Expected HashedPassword %s, got %s", expectedHashedPassword, output.HashedPassword)
	}
}
