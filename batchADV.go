// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
)

// BatchADV holds the Batch Header and Batch Control and all Entry Records for ADV Entries
type BatchADV struct {
	Batch
}

// NewBatchADV returns a *BatchADV
func NewBatchADV(bh *BatchHeader) *BatchADV {
	batch := new(BatchADV)
	batch.SetADVControl(NewADVBatchControl())
	batch.SetHeader(bh)
	return batch
}

// Validate checks valid NACHA batch rules. Assumes properly parsed records.
func (batch *BatchADV) Validate() error {

	if batch.Header.StandardEntryClassCode != "ADV" {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.StandardEntryClassCode, "ADV")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "StandardEntryClassCode", Msg: msg}
	}
	if batch.Header.ServiceClassCode != 280 {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.ServiceClassCode, "ADV")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "ServiceClassCode", Msg: msg}
	}
	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}
	// Add configuration and type specific validation for this type.
	for _, entry := range batch.ADVEntries {

		if entry.Category == CategoryForward {
			switch entry.TransactionCode {
			case 81, 82, 83, 84, 85, 86, 87, 88:
			default:
				msg := fmt.Sprintf(msgBatchTransactionCode, entry.TransactionCode, "ADV")
				return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "TransactionCode", Msg: msg}
			}
		}
		// ToDo:  If NOC and Returns for ADV check Addenda98 and Addenda99
	}
	return nil
}

// Create takes Batch Header and Entries and builds a valid batch
func (batch *BatchADV) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}
	// Additional steps specific to batch type
	// ...
	return batch.Validate()
}
