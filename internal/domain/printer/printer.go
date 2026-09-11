package printer

type ConnectionType string
type PrinterStatus string
type PrinterProtocol string

const (
	ConnectionUSB     ConnectionType = "usb"
	ConnectionSerial  ConnectionType = "serial"
	ConnectionNetwork ConnectionType = "network"
	ConnectionSystem  ConnectionType = "system"
	ConnectionUnknown ConnectionType = "unknown"

	ProtocolEscpos  PrinterProtocol = "escpos"
	ProtocolCups    PrinterProtocol = "cups"
	ProtocolWindows PrinterProtocol = "windows"
	ProtocolRawTCP  PrinterProtocol = "raw_tcp"
	ProtocolUnknown PrinterProtocol = "unknown"

	StatusUnknown  PrinterStatus = "unknown"
	StatusReady    PrinterStatus = "ready"
	StatusPrinting PrinterStatus = "printing"
	StatusOffline  PrinterStatus = "offline"
	StatusPaperOut PrinterStatus = "paper_out"
	StatusError    PrinterStatus = "error"
	StatusBusy     PrinterStatus = "busy"
)

type Printer struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Manufacturer string          `json:"manufacturer,omitempty"`
	Model        string          `json:"model,omitempty"`
	Connection   ConnectionType  `json:"connection"`
	Address      string          `json:"address,omitempty"`
	Status       PrinterStatus   `json:"status"`
	Protocol     PrinterProtocol `json:"protocol"`
	IsDefault    bool            `json:"is_default"`
}
