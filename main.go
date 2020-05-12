package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	MinimumHeaderSize = 12
	MaximumeaderSize  = 14
)

func (r *DataRecord) GetMessageLength() (int64, error) {
	switch r.Header.LocalMessageType {
	case 0:
		return 1, nil
	}
	return 0, nil
}

func (t *DataRecordMessageType) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, errors.New("data record message type must be defined")
	}

	switch *t {
	case DataRecordMessageType_Definition:
		return json.Marshal("DEFINITION")
	case DataRecordMessageType_Data:
		return json.Marshal("DATA")
	default:
		return nil, fmt.Errorf("unknown data record message type %d", *t)
	}
}

func (t *DataRecordHeaderType) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, errors.New("data record header type must be defined")
	}

	switch *t {
	case DataRecordHeaderType_Normal:
		return json.Marshal("NORMAL")
	case DataRecordHeaderType_CompressedTimestamp:
		return json.Marshal("COMPRESSED_TIMESTAMP")
	default:
		return nil, fmt.Errorf("unknown data record header type %d", *t)
	}
}

func (h *DataRecordHeader) Unmarshal(data []byte) error {
	if data == nil || len(data) != 1 {
		return errors.New("a data record header must be exactly 1 byte")
	}
	exploded := ExplodeByte(uint8(data[0]))
	if len(exploded) != 8 {
		return errors.New("bit represention of a byte must be of length 8")
	}

	isNormalHeader, err := strconv.ParseInt(string(exploded[0:1]), 2, 8)
	if err != nil {
		return err
	}
	if isNormalHeader == 0 {
		h.Type = DataRecordHeaderType_Normal
		localMessageType, err := strconv.ParseInt(string(exploded[4:8]), 2, 8)
		if err != nil {
			return err
		}
		h.LocalMessageType = uint8(localMessageType)
		messageType, err := strconv.ParseInt(string(exploded[1:2]), 2, 8)
		if err != nil {
			return err
		}
		if messageType == 1 {
			h.MessageType = DataRecordMessageType_Definition
		} else if messageType == 0 {
			h.MessageType = DataRecordMessageType_Data
		}
		developerData, err := strconv.ParseInt(string(exploded[2:3]), 2, 8)
		if err != nil {
			return err
		}
		if developerData == 1 {
			h.DeveloperData = true
		}
	} else if isNormalHeader == 1 {
		h.Type = DataRecordHeaderType_CompressedTimestamp
		h.MessageType = DataRecordMessageType_Data
		timeOffset, err := strconv.ParseInt(string(exploded[3:8]), 2, 8)
		if err != nil {
			return err
		}
		h.TimeOffset = uint8(timeOffset)
		localMessageType, err := strconv.ParseInt(string(exploded[1:3]), 2, 8)
		if err != nil {
			return err
		}
		h.LocalMessageType = uint8(localMessageType)
	}

	return nil
}

func (m *DefinitionMessage) DataMessageSize() int64 {
	size := int64(0)

	for _, feildDefinition := range m.VariableContent.Fields {
		size += int64(feildDefinition.Size)
	}

	return size
}

func (c *DefinitionMessageFixedContent) Unmarshal(data []byte) error {
	if data == nil || len(data) != 5 {
		return errors.New("a definition message fixed content must be exactly 5 bytes")
	}
	c.Architecture = data[1]
	c.GlobalMessageNumber = binary.BigEndian.Uint16(data[2:4])
	c.NumFields = data[4]
	return nil
}

func (c *DefinitionMessageVariableContent) Unmarshal(data []byte) error {
	if c == nil || data == nil || len(data)%3 != 0 {
		return errors.New("a definition message variable content must be a multiple of 3")
	}
	c.Fields = make([]FieldDefinition, len(data)/3)

	for i := 0; i < len(data); i += 3 {
		baseTypeExploded := ExplodeByte(uint8(data[i+2]))

		baseTypeNumber, err := strconv.ParseInt(string(baseTypeExploded[3:8]), 2, 8)
		if err != nil {
			return err
		}

		endianAbility, err := strconv.ParseInt(string(baseTypeExploded[0:1]), 2, 8)
		if err != nil {
			return err
		}

		c.Fields[i/3] = FieldDefinition{
			Number: uint8(data[i]),
			Size:   uint8(data[i+1]),
			BaseType: BaseType{
				Number:        uint8(baseTypeNumber),
				EndianAbility: uint8(endianAbility),
			},
		}
	}

	return nil
}

func (m *DataMessage) Unmarshal(def *DefinitionMessage, data []byte) error {
	if m == nil || def == nil || int64(len(data)) != def.DataMessageSize() {
		return errors.New("data message size incorrect") // TODO: better error message
	}

	m.Fields = make([]uint64, len(def.VariableContent.Fields))
	offset := 0
	for i := 0; i < len(m.Fields); i++ {
		switch def.VariableContent.Fields[i].BaseType.Number {
		case 0, 1, 2, 7, 10, 13:
			m.Fields[i] = uint64(data[offset])
		case 3, 4, 11:
			m.Fields[i] = uint64(binary.BigEndian.Uint16(data[offset : offset+int(def.VariableContent.Fields[i].Size)]))
		case 5, 6, 8, 12:
			m.Fields[i] = uint64(binary.BigEndian.Uint32(data[offset : offset+int(def.VariableContent.Fields[i].Size)]))
		case 9, 14, 15, 16:
			m.Fields[i] = binary.BigEndian.Uint64(data[offset : offset+int(def.VariableContent.Fields[i].Size)])
		default:
			return errors.New("unknown base type number")
		}
		offset += int(def.VariableContent.Fields[i].Size)
	}

	return nil
}

// Unmarshal populates a Header struct from raw binary data
func (h *FileHeader) Unmarshal(data []byte) error {
	if data == nil || len(data) < 1 {
		return errors.New("data must be defined")
	}

	h.Size = data[0]
	if h.Size != MinimumHeaderSize && h.Size != MaximumeaderSize {
		return fmt.Errorf("valid header sizes are %d and %d", MinimumHeaderSize, MaximumeaderSize)
	}
	if uint8(len(data)) != h.Size {
		return fmt.Errorf("data must be at least of size %d", h.Size)
	}

	h.ProtocolVersion = data[1]
	h.ProfileVersion = binary.BigEndian.Uint16(data[2:4])
	h.DataSize = binary.BigEndian.Uint32(data[4:8])
	h.DataType = string(data[8:12])
	if h.Size == MaximumeaderSize {
		h.CRC = binary.LittleEndian.Uint16(data[12:14])
	}
	return nil
}

// ExplodeByte converts a byte into its binary string representation
func ExplodeByte(data uint8) string {
	var b strings.Builder

	for i := 7; i >= 0; i-- {
		if (data & (1 << uint(i))) != 0 {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}

	return b.String()
}

func (dm *DataMessage) ReadAndUnmarshal(ctx context.Context, r io.Reader) (int, error) {
	if dm == nil {
		return 0, errors.New("definition message is nil")
	}

	def := ctx.Value("CURRENT_DEFINITION_RECORD").(DataRecord)
	if def.DefinitionMessage == nil || def.DefinitionMessage.VariableContent == nil {
		return 0, errors.New("definition message must preceed a data message")
	}

	dm.Fields = []uint64{}

	var totalBytesRead int

	buf := make([]byte, def.DefinitionMessage.DataMessageSize())
	n, err := r.Read(buf)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	if err := dm.Unmarshal(def.DefinitionMessage, buf); err != nil {
		return totalBytesRead, err
	}

	return totalBytesRead, nil
}

func (dm *DefinitionMessage) ReadAndUnmarshal(ctx context.Context, r io.Reader) (int, error) {
	if dm == nil {
		return 0, errors.New("definition message is nil")
	}

	var totalBytesRead int

	dm.FixedContent = new(DefinitionMessageFixedContent)
	dm.VariableContent = new(DefinitionMessageVariableContent)

	buf := make([]byte, 5)
	n, err := r.Read(buf)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	if err := dm.FixedContent.Unmarshal(buf); err != nil {
		return totalBytesRead, err
	}

	buf = make([]byte, int64(dm.FixedContent.NumFields)*3)
	n, err = r.Read(buf)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	if err := dm.VariableContent.Unmarshal(buf); err != nil {
		return totalBytesRead, err
	}

	if ctx.Value("DATA_RECORD_HEADER").(*DataRecordHeader).DeveloperData {
		buf = make([]byte, 1)
		n, err = r.Read(buf)
		totalBytesRead += n
		if err != nil {
			return totalBytesRead, err
		}
		dm.FixedContent.NumFields += uint8(buf[0])

		vc := new(DefinitionMessageVariableContent)
		buf = make([]byte, int64(buf[0])*3)
		n, err = r.Read(buf)
		totalBytesRead += n
		if err != nil {
			return totalBytesRead, err
		}

		if err := vc.Unmarshal(buf); err != nil {
			return totalBytesRead, err
		}
		dm.VariableContent.Fields = append(dm.VariableContent.Fields, vc.Fields...)
	}

	return totalBytesRead, nil
}

func (dr *DataRecord) ReadAndUnmarshal(ctx context.Context, r io.Reader) (int, error) {
	if dr == nil {
		return 0, errors.New("data record is nil")
	}

	var totalBytesRead int

	dr.Header = new(DataRecordHeader)

	// no matter what the message type, the header size is the same
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	if err := dr.Header.Unmarshal(buf); err != nil {
		return totalBytesRead, err
	}
	switch dr.Header.MessageType {
	case DataRecordMessageType_Definition:
		dr.DefinitionMessage = new(DefinitionMessage)
		n, err := dr.DefinitionMessage.ReadAndUnmarshal(context.WithValue(context.Background(), "DATA_RECORD_HEADER", (*dr).Header), r)
		totalBytesRead += n
		if err != nil {
			return totalBytesRead, err
		}
	case DataRecordMessageType_Data:
		dr.DataMessage = new(DataMessage)
		n, err := dr.DataMessage.ReadAndUnmarshal(ctx, r)
		totalBytesRead += n
		if err != nil {
			return totalBytesRead, err
		}
	default:
		return totalBytesRead, errors.New("unkown message type")
	}

	return totalBytesRead, nil
}

func (h *FileHeader) ReadAndUnmarshal(r io.Reader) (int, error) {
	buf := make([]byte, 14)
	n, err := r.Read(buf)
	if err != nil {
		return n, err
	}
	return n, h.Unmarshal(buf)
}

func (f *File) ReadAndUnmarshal(r io.Reader) (int, error) {
	if f == nil {
		return 0, errors.New("file is nil")
	}

	if f.Header == nil {
		f.Header = new(FileHeader)
	}

	if f.DataRecords == nil {
		f.DataRecords = []DataRecord{}
	}

	var totalBytesRead int

	n, err := f.Header.ReadAndUnmarshal(r)
	totalBytesRead += n
	if err != nil {
		return totalBytesRead, err
	}

	recordsProcessed := 0

	var currentDefinitionMessage DataRecord

	dataRecordsBytesLeftToProcess := int(f.Header.DataSize)
	for dataRecordsBytesLeftToProcess > 0 {
		dr := new(DataRecord)

		ctx := context.Background()
		if recordsProcessed > 0 {
			ctx = context.WithValue(ctx, "CURRENT_DEFINITION_RECORD", currentDefinitionMessage)
		}

		n, err := dr.ReadAndUnmarshal(ctx, r)
		totalBytesRead += n
		dataRecordsBytesLeftToProcess -= n
		if err != nil {
			return totalBytesRead, err
		}

		if dr.Header.MessageType == DataRecordMessageType_Definition {
			currentDefinitionMessage = *dr
		}

		f.DataRecords = append(f.DataRecords, *dr)
		recordsProcessed++
	}
	return totalBytesRead, nil
}

func main() {
	rawFile, err := os.Open(os.Getenv("FILE_LOCATION"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer rawFile.Close()

	fitFile := new(File)
	if _, err := fitFile.ReadAndUnmarshal(rawFile); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	data, err := json.Marshal(fitFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(string(data))
}
