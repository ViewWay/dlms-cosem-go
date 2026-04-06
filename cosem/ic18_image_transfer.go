package cosem

import (
	"github.com/ViewWay/dlms-cosem-go/core"
)

// ImageTransferStatus represents the image transfer status
type ImageTransferStatus uint8

const (
	ImageTransferStatusIdle                    ImageTransferStatus = 0
	ImageTransferStatusInitiated               ImageTransferStatus = 1
	ImageTransferStatusInitiatedForVerifying   ImageTransferStatus = 2
	ImageTransferStatusVerifyingInitiated      ImageTransferStatus = 3
	ImageTransferStatusVerificationFailed      ImageTransferStatus = 4
	ImageTransferStatusVerificationSuccessful  ImageTransferStatus = 5
	ImageTransferStatusImageActivated          ImageTransferStatus = 6
	ImageTransferStatusImageNotActivated       ImageTransferStatus = 7
)

// ImageTransfer is the COSEM Image Transfer interface class (IC 35).
type ImageTransfer struct {
	LogicalName     core.ObisCode
	ImageBlockSize  uint16
	ImageFirstBlock []byte
	ImageBlockCount uint32
	ImageReference  []byte
	ImageIdent      uint8
	TransferStatus  ImageTransferStatus
	Version         uint8
}

func (it *ImageTransfer) ClassID() uint16               { return core.ClassIDImageTransfer }
func (it *ImageTransfer) GetLogicalName() core.ObisCode { return it.LogicalName }
func (it *ImageTransfer) GetVersion() uint8             { return it.Version }
func (it *ImageTransfer) Access(attr int) core.CosemAttributeAccess {
	switch attr {
	case 2, 3, 4, 5, 6, 7:
		return core.CosemAttributeAccess{Access: core.AccessRead}
	default:
		return core.CosemAttributeAccess{Access: core.AccessNone}
	}
}

func (it *ImageTransfer) MarshalBinary() ([]byte, error) {
	return core.EnumData(it.TransferStatus).ToBytes(), nil
}

func (it *ImageTransfer) UnmarshalBinary(data []byte) error {
	elem, _, err := core.DlmsDataFromBytes(data)
	if err != nil {
		return err
	}
	if v, ok := elem.(core.EnumData); ok {
		it.TransferStatus = ImageTransferStatus(v)
	}
	return nil
}

func (it *ImageTransfer) InitiateImage(blockSize uint16, blockCount uint32) {
	it.ImageBlockSize = blockSize
	it.ImageBlockCount = blockCount
	it.TransferStatus = ImageTransferStatusInitiated
}

func (it *ImageTransfer) VerifyImage() bool {
	it.TransferStatus = ImageTransferStatusVerificationSuccessful
	return true
}

func (it *ImageTransfer) ActivateImage() error {
	it.TransferStatus = ImageTransferStatusImageActivated
	return nil
}
