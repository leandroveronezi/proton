package proton

type flavor int

const (
	Chrome flavor = 1
	Edge   flavor = 2
)

type Config struct {
	Title           string
	Url             string
	Debug           bool
	UserDataDir     string
	UserDataDirKeep bool
	Width           int
	Height          int
	WindowState     WindowState
	Flavor          flavor
	Args            []string
	BrowserBinary   string
}

var DefaultBrowserArgs = []string{
	"--disable-background-networking",
	"--disable-background-timer-throttling",
	"--disable-backgrounding-occluded-windows",
	"--disable-breakpad",
	"--disable-client-side-phishing-detection",
	"--disable-default-apps",
	"--disable-dev-shm-usage",
	"--disable-infobars",
	"--disable-extensions",
	"--disable-features=site-per-process",
	"--disable-hang-monitor",
	"--disable-ipc-flooding-protection",
	"--disable-popup-blocking",
	"--disable-prompt-on-repost",
	"--disable-renderer-backgrounding",
	"--disable-sync",
	"--disable-translate",
	"--metrics-recording-only",
	"--no-first-run",
	"--safebrowsing-disable-auto-update",
	"--enable-automation",
	"--password-store=basic",
	"--use-mock-keychain",
	"--disable-dinosaur-easter-egg",
	"--disable-windows10-custom-titlebar",
	"--no-default-browser-check",
}
