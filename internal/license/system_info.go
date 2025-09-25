package license

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

func GetMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if len(iface.HardwareAddr) >= 6 {
				return strings.ToUpper(hex.EncodeToString(iface.HardwareAddr))
			}
		}
	}

	return ""
}

func GetIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}

	defer func() { _ = conn.Close() }()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func GetSystemFingerprint() string {
	hostname, _ := os.Hostname()
	osInfo := runtime.GOOS + "/" + runtime.GOARCH

	data := fmt.Sprintf("%s:%s:%s:%s", hostname, osInfo, GetMACAddress(), GetIPAddress())

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func GetMultiHostFingerprint() string {
	hostname, _ := os.Hostname()
	osInfo := runtime.GOOS + "/" + runtime.GOARCH

	data := fmt.Sprintf("%s:%s", hostname, osInfo)
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
