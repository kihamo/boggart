package pantum

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/pantum/client"
	"github.com/kihamo/boggart/providers/pantum/models"
)

type ProductID uint64
type CartridgeStatus int64
type DrumStatus int64

const (
	ProductIDUnknown        ProductID = 0
	ProductIDP3012NET       ProductID = 0x1830
	ProductIDP3100DLP3255DN ProductID = 0x6830
	ProductIDP3010D         ProductID = 0x0EC4
	ProductIDP3010DW        ProductID = 0x0EC5
	ProductIDP3060DW        ProductID = 0x0EC6
	ProductIDP3300          ProductID = 0x0EC7
	ProductIDP3300DN        ProductID = 0x0EC8
	ProductIDP3300DW        ProductID = 0x0EC9
	ProductIDM6700D         ProductID = 0x0ECA
	ProductIDM6700DN        ProductID = 0x0EE1
	ProductIDM6760D         ProductID = 0x0F02
	ProductIDM6700DW        ProductID = 0x0ECB
	ProductIDM6760DW        ProductID = 0x0ECC
	ProductIDM7100D         ProductID = 0x0EE2
	ProductIDM7100DN        ProductID = 0x0ECF
	ProductIDM7100DW        ProductID = 0x0ED0
	ProductIDM7160DW        ProductID = 0x0F03
	ProductIDM6800FDW       ProductID = 0x0ECD
	ProductIDM6860FDW       ProductID = 0x0ECE
	ProductIDM6860FDN       ProductID = 0x0EF7
	ProductIDM7200FD        ProductID = 0x0ED1
	ProductIDM7200FDN       ProductID = 0x0ED2
	ProductIDM7200FDW       ProductID = 0x0ED3

	CartridgeStatusUnknown CartridgeStatus = -1
	CartridgeStatusNormal  CartridgeStatus = 0
	CartridgeNotDetected   CartridgeStatus = 1
	CartridgeMismatch      CartridgeStatus = 2
	CartridgeExpired       CartridgeStatus = 3
	CartridgeTonerLow      CartridgeStatus = 4

	DrumStatusUnknown     DrumStatus = -1
	DrumStatusNormal      DrumStatus = 0
	DrumStatusUninstalled DrumStatus = 1
	DrumStatusUnmatched   DrumStatus = 2
	DrumStatusInvalid     DrumStatus = 3
	DrumStatusExpired     DrumStatus = 4
)

var systemStatus = map[int64]string{
	-1: "Unknown Status",
	0:  "Initialization",
	1:  "Initialization",
	2:  "Ready",
	3:  "Ready",
	4:  "Sleep",
	5:  "Printing",
	10: "Warming Up",
	11: "Printing",
	14: "Manual feeder tray is out of paper",
	15: "Manual Feeder Tray Feed failure",
	16: "Manual Feeder Tray Feed Jam",
	17: "Manual Feeder Tray Paper Source Mismatch",
	18: "Manual Feeder Tray Paper Setting Mismatch",
	19: "Automatic feeder tray is out of paper",
	20: "Automatic Feeder Tray Feed failure",
	21: "Automatic Feeder Tray Feed Jam",
	22: "Automatic Feeder Tray Paper Source Mismatch",
	23: "Automatic Feeder Tray Paper Setting Mismatch",
	24: "Middle Jam",
	25: "Exit jam",
	26: "Paper Jam in the Duplex Unit",
	27: "Middle Jam Not Clear",
	28: "Exit Jam Not Clear",
	29: "Paper Jam in the Duplex Unit",
	30: "Front cover open",
	31: "Rear cover open",
	32: "Printer internal failure 01",
	33: "Printer internal failure 02",
	34: "Printer internal failure 03",
	35: "Printer internal failure 04",
	36: "Printer internal failure 05",
	37: "Printer internal failure 06",
	38: "Printer internal failure 07",
	39: "Printer internal failure 08",
	40: "Printer internal failure 09",
	41: "Printer internal failure 10",
	42: "Printer internal failure 11",
	43: "Printer internal failure 12",
	44: "Printer internal failure 13",
	45: "Canceling",
	46: "Cartridge not detected",
	47: "Cartridge mismatch",
	48: "Cartridge life has expired",
	49: "Drum unit is uninstalled",
	50: "Type of the drum unit is unmatched",
	51: "Life of the drum unit will come to an end",
	52: "Drum unit is invalid",
	53: "Toner Low / Life of the drum unit will come to an end",
	54: "Automatic feeder tray is out of paper + Manual feeder tray is out of paper",
}

var scanStatus = map[int64]string{
	1:  "Initialization",
	2:  "Ready",
	3:  "Processing",
	4:  "Scanning",
	5:  "Scanning",
	6:  "Canceling",
	7:  "Scan next page: start",
	8:  "Ready",
	9:  "Ready",
	10: "ADF Paper empty",
	11: "ADF Paper jam",
	12: "ADF Paper jam",
	13: "Close ADF upper cover",
	14: "Low memory",
	15: "The mail size exceeds the limit",
	16: "Communication error 21",
	17: "Communication error 22",
	18: "Communication error 23",
	19: "Communication error 24",
	20: "Communication error 25",
	21: "Communication error 26",
	22: "Communication error 27",
	23: "Fail to get scanned file",
	24: "Username or Password Error",
	25: "File size error",
	26: "fail to send scanned file",
	27: "Scanner internal error 11",
	28: "Scanner internal error 12",
	29: "Scanner internal error 13",
	30: "Scanner internal error 14",
	31: "Scanner internal error 15",
	32: "Scanner internal error 16",
}

var copyStatus = map[int64]string{
	1: "Ready",
	2: "Copying",
	3: "Canceling",
	4: "Back side：Ok",
	5: "Back side of ID card：Start",
	6: "Low memory",
}

var faxStatus = map[int64]string{
	1:  "Ready",
	2:  "Dialing",
	3:  "Connecting",
	4:  "Sending",
	6:  "Sending",
	7:  "Recieving",
	8:  "Recieving",
	9:  "Scanning",
	10: "Printing",
	11: "Canceling",
	12: "Canceled!",
	13: "Fax done!",
	14: "Received!",
	15: "Scan next page: start",
	19: "Comm. Failed",
	20: "Communication failed + Wait for redial",
	21: "No reply",
	22: "No answer + Wait for redial",
	26: "Receive memory is full",
	27: "Send memory is full",
	28: "Dialing",
	31: "Canceling junk",
	32: "Delay send saved",
	33: "Fax done!",
	34: "Incoming call",
	35: "Phone",
	36: "Tap Start recei.",
	37: "Phone",
	38: "Tap Start recei.",
	39: "No answer + Please check your fax line",
	41: "Communication failed + Please check your fax line",
	43: "Initialization",
}

func isSeparator(r rune) bool {
	return r == '(' || r == ')' || r == ','
}

func PrinterStatus(errorFlag, modules string) string {
	errorCode, err := strconv.ParseInt(errorFlag, 10, 64)
	if err != nil {
		return systemStatus[-1]
	}

	// check system status
	if s, ok := systemStatus[errorCode]; ok {
		return s
	}

	var (
		offsetScan int64 = 55 + 1
		offsetCopy int64 = 94 + 1
		offsetFax  int64 = 101 + 1
	)

	parts := strings.Split(modules, "-")
	if len(parts) == 3 {
		offsetScan, _ = strconv.ParseInt(parts[0], 10, 64)
		offsetCopy, _ = strconv.ParseInt(parts[1], 10, 64)
		offsetFax, _ = strconv.ParseInt(parts[2], 10, 64)
	}

	// check scan status
	if s, ok := scanStatus[errorCode-offsetScan]; ok {
		return s
	}

	// check copy status
	if s, ok := copyStatus[errorCode-offsetCopy]; ok {
		return s
	}

	// check fax status
	if s, ok := faxStatus[errorCode-offsetFax]; ok {
		return s
	}

	return systemStatus[-1]
}

func ProductIDConvert(productID string) ProductID {
	id, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		return ProductIDUnknown
	}

	ret := ProductID(id)
	if !ret.IsAProductID() {
		return ProductIDUnknown
	}

	return ret
}

func CartridgeStatusConvert(status string) CartridgeStatus {
	id, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		return CartridgeStatusUnknown
	}

	ret := CartridgeStatus(id)
	if !ret.IsACartridgeStatus() {
		return CartridgeStatusUnknown
	}

	return ret
}

func DrumStatusConvert(status string) DrumStatus {
	id, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		return DrumStatusUnknown
	}

	ret := DrumStatus(id)
	if !ret.IsADrumStatus() {
		return DrumStatusUnknown
	}

	return ret
}

func DatabaseToMap(properties models.Properties) map[string]interface{} {
	db := make(map[string]interface{}, len(properties))

	for _, p := range properties {
		switch p.Type {
		case "InputCheckbox":
			db[p.Name], _ = strconv.ParseBool(p.Value)
		case "InputIpaddr":
			db[p.Name] = net.ParseIP(p.Value)
		default: // OnlyValue, StaticValue, InputText, InputPassword, Selection, InputShort, InputRadio
			db[p.Name] = p.Value
		}
	}

	return db
}

type property struct {
	Value  string `json:"value"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Module string `json:"module"`
	Index  int    `json:"index"`
}

type Client struct {
	*client.Pantum
}

func New(address *url.URL, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().
		WithSchemes([]string{address.Scheme}).
		WithHost(address.Host)

	cl := &Client{
		Pantum: client.NewHTTPClientWithConfig(nil, cfg),
	}

	if rt, ok := cl.Pantum.Transport.(*httptransport.Runtime); ok {
		rt.Consumers["text/html"] = omConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func omConsumer() runtime.Consumer {
	jsonConsumer := runtime.JSONConsumer().Consume

	return runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		jsonBody := bytes.NewBuffer(nil)

		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanLines)

		prop := &property{}

		for scanner.Scan() {
			parts := bytes.FieldsFunc(scanner.Bytes(), isSeparator)

			if len(parts) > 4 {
				for i, v := range parts {
					parts[i] = bytes.Trim(v, " '\"")
				}

				valueDecode, _ := base64.StdEncoding.DecodeString(string(parts[1]))

				prop.Value = string(valueDecode)
				prop.Name = string(parts[2])
				prop.Type = strings.TrimPrefix(string(parts[3]), "SN.TYPE.")
				prop.Module = strings.TrimPrefix(string(parts[4]), "MODULE_")
				prop.Index, _ = strconv.Atoi(string(parts[5]))

				if line, err := json.Marshal(prop); err == nil {
					if jsonBody.Len() == 0 {
						jsonBody.WriteRune('[')
					} else {
						jsonBody.WriteByte(',')
					}

					jsonBody.Write(line)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		if jsonBody.Len() > 0 {
			jsonBody.WriteByte(']')
		}

		return jsonConsumer(jsonBody, data)
	})
}
