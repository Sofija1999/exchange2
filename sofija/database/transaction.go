package database

import (
	"github.com/pkg/errors"
)

func TxHandler(tx *Tx, err error, recovery interface{}) error {
	if recovery != nil {
		switch v := recovery.(type) {
		case error:
			err = v
		case string:
			err = errors.New(v)
		default:
			err = errors.New("recovered from panic")
		}
	}

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(err, rollbackErr.Error())
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}
