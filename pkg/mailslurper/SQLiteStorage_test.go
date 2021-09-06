package mailslurper

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSQLLiteStorage(t *testing.T) *SQLiteStorage {
	t.Helper()

	logger := logrus.NewEntry(logrus.New())
	storage := NewSQLiteStorage(&ConnectionInformation{
		Filename: ":memory:",
	}, logger)
	require.NoError(t, storage.Connect())
	require.NoError(t, storage.Create())
	return storage
}

func TestSQLiteStorage_GetMailCollection(t *testing.T) {
	t.Run("should list mailitems", func(t *testing.T) {
		wantCount := 1

		storage := newSQLLiteStorage(t)
		_, err := storage.StoreMail(&MailItem{
			ID:          "id1",
			FromAddress: "from@email.com",
			ToAddresses: MailAddressCollection{
				"to1@email.com",
				"to2@email.com",
			},
			Subject:     "Subject value",
			Body:        "body value",
			ContentType: "text/plain",
			TextBody:    "text body value",
			HTMLBody:    "html body value",
		})
		require.NoError(t, err)
		mailSearch := &MailSearch{}
		emailCount, err := storage.GetMailCount(mailSearch)
		require.NoError(t, err)
		assert.Equal(t, wantCount, emailCount)
		emails, err := storage.GetMailCollection(0, 100, mailSearch)
		require.NoError(t, err)
		require.Len(t, emails, 1)
	})
}
