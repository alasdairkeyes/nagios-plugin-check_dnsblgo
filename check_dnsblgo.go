package main

    // Packages to import
import (
    "flag"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
    "os"
    "net"
)

    // Declare vars
var version = "0.1.0"

// Messages to exit codes mapping
var statusmap = map[string]int{
    "OK":       0,
    "WARNING":  1,
    "CRITICAL": 2,
}

    // Get the directory of this binary
func getCwd() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        exit("CRITICAL", "Failed to get directory of tool: " + err.Error())
    }
    return dir
}

    // Get args from CLI
func getArgs() (string, string){

        // Define vars and defaults
    var ipString string
    var blocklistFilename string
    var help bool

    blocklistsDefaultFilename := "blocklists.txt"
    blocklistsDefaultPath := filepath.Join(getCwd(), blocklistsDefaultFilename)

        // Declare the flags
    flag.StringVar(
        &ipString,
        "ip",
        "",
        "IPv4 address to check")
    flag.StringVar(
        &blocklistFilename,
        "file",
        blocklistsDefaultPath,
        "Filename with list of blocklists to check, one per line")
    flag.BoolVar(
        &help,
        "help",
        false,
        "Display the help information")

        // Parse
    flag.Parse()

        // Show help message and exit if required
    if (help == true) {
        usage()
    }

        // Validate that we have an IPv4 address
    ip := net.ParseIP(ipString)
    if ip == nil || ip.To4() == nil {
        exit("CRITICAL", "Invalid IPv4 address " + ipString);
    }

    return ipString, blocklistFilename
}

    // Helpscreen showing usage
func usage() {
    flag.PrintDefaults()
    fmt.Println("\ncheck_dnsbl uses the system's default DNS resolver")
    fmt.Println("Version " + version)
    os.Exit(0)
}


func output(ipString string, blocklists []string, hits []string, lookupErrors []string) {
    // Output critical for hits
    if len(hits) > 0 {
        exit(
            "CRITICAL",
            fmt.Sprintf("%s is listed on %s", ipString, strings.Join(hits, ", ")))
    }

        // Display that there were lookup errors
    if len(lookupErrors) > 0 {
        exit(
            "WARNING",
            fmt.Sprintf("DNS lookup errors on %s blocklists", strings.Join(lookupErrors, ", ")))
    }

        // All OK
    exit(
        "OK",
        fmt.Sprintf("%s is not listed on %d blocklists", ipString, len(blocklists)));
}


func main() {

        // Get args
    ipString, blocklistFilename := getArgs()

        // Load list
    blocklists := importBlocklists(blocklistFilename);

        // DNS checks
    var hits, lookupErrors = checkIP(ipString, blocklists)

        // Display the results
    output(ipString, blocklists, hits, lookupErrors)
}

    // Reverse an IP
func reverseIP(ipString string) string {
    ipParts := strings.Split(ipString, ".")
    sort.Sort(sort.Reverse(sort.StringSlice(ipParts)))
    return strings.Join(ipParts, ".")
}

    // Check an IP against a list of blocklists
func checkIP(ipString string, blocklists []string) ([]string, []string) {

    var hits []string;
    var lookupErrors []string

        // Define regexp for detecting if no host
    noHostRegexp, _ := regexp.Compile(`no such host$`);

    reverseIPString := reverseIP(ipString)

    for _, blocklist := range blocklists {
        blocklistIPHostname := getBlocklistIPHostname(reverseIPString, blocklist)

            // Perform the lookup
        _, err := net.LookupHost(blocklistIPHostname);

        if err != nil {
            // This is less than pleasant, no specific errors seem to be returned
            if ! noHostRegexp.MatchString(err.Error()) {
                lookupErrors = append(lookupErrors, blocklist)
            }

        } else {
            // The host is on blocklist
            hits = append(hits, blocklist);
        }
    }

    return hits, lookupErrors
}

    // Join reversed IP with given hostname of blocklist
func getBlocklistIPHostname (reverseIPString, blocklist string) string {
    return reverseIPString + "." + blocklist
}

    // Load the contents of file and extract the blocklists to use
func importBlocklists(filename string) []string {
        // Load file or report error
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        exit("CRITICAL", fmt.Sprintf("%s", filename, err));
    }
    contentString := fmt.Sprintf("%s", content);

        // Define regexp to remove unneeded lines
    ignoreLineRegexp, _ := regexp.Compile(`^\s*(#.*)?$`);

        // Split file contents into array based on newlines
    blocklists := strings.Split(contentString, "\n");

        // New array to hold valid blocklists
    var cleanBlocklists []string

        // Iterate blocklist array
    for _, value := range blocklists {

        // Clean whitespace
        value = strings.Trim(value, " ");

        // Copy valid lines into new array
        if ! ignoreLineRegexp.MatchString(value) {
            cleanBlocklists = append(cleanBlocklists, value)
        }
    }

    return cleanBlocklists
}

    // Format output to Nagios expected format
func exit (status, message string) {

    fmt.Printf("%s: %s", status, message);
    os.Exit(statusmap[status])
}
