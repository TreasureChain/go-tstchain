package mru

import (
	"github.com/TreasureChain/go-tstchain/swarm/storage"
)

// updateHeader models the non-payload components of a Resource Update
type updateHeader struct {
	UpdateLookup        // UpdateLookup contains the information required to locate this resource (components of the search key used to find it)
	multihash    bool   // Whether the data in this Resource Update should be interpreted as multihash
	metaHash     []byte // SHA3 hash of the metadata chunk (less ownerAddr). Used to prove ownerhsip of the resource.
}

const metaHashLength = storage.KeyLength

// updateLookupLength bytes
// 1 byte flags (multihash bool for now)
// 32 bytes metaHash
const updateHeaderLength = updateLookupLength + 1 + metaHashLength

// binaryPut serializes the resource header information into the given slice
func (h *updateHeader) binaryPut(serializedData []byte) error {
	if len(serializedData) != updateHeaderLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to serialize updateHeaderLength. Expected %d, got %d", updateHeaderLength, len(serializedData))
	}
	if len(h.metaHash) != metaHashLength {
		return NewError(ErrInvalidValue, "updateHeader.binaryPut called without metaHash set")
	}
	if err := h.UpdateLookup.binaryPut(serializedData[:updateLookupLength]); err != nil {
		return err
	}
	cursor := updateLookupLength
	copy(serializedData[cursor:], h.metaHash[:metaHashLength])
	cursor += metaHashLength

	var flags byte
	if h.multihash {
		flags |= 0x01
	}

	serializedData[cursor] = flags
	cursor++

	return nil
}

// binaryLength returns the expected size of this structure when serialized
func (h *updateHeader) binaryLength() int {
	return updateHeaderLength
}

// binaryGet restores the current updateHeader instance from the information contained in the passed slice
func (h *updateHeader) binaryGet(serializedData []byte) error {
	if len(serializedData) != updateHeaderLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to read updateHeaderLength. Expected %d, got %d", updateHeaderLength, len(serializedData))
	}

	if err := h.UpdateLookup.binaryGet(serializedData[:updateLookupLength]); err != nil {
		return err
	}
	cursor := updateLookupLength
	h.metaHash = make([]byte, metaHashLength)
	copy(h.metaHash[:storage.KeyLength], serializedData[cursor:cursor+storage.KeyLength])
	cursor += metaHashLength

	flags := serializedData[cursor]
	cursor++

	h.multihash = flags&0x01 != 0

	return nil
}
