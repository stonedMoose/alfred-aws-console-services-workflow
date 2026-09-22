package awsconfig

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

// Profile is a named AWS profile from the shared credentials or config file.
type Profile struct {
	Name   string
	Region string
}

// configProfilePrefix is what the config file puts before every profile name
// but the default one.
const configProfilePrefix = "profile "

var awsProfiles []Profile

// GetAwsProfiles lists the profiles of the credentials and config files,
// reading them once.
func GetAwsProfiles() []Profile {
	if len(awsProfiles) <= 0 {
		loadAwsProfiles()
	}
	return awsProfiles
}

// FindProfile returns the profile of that name, or nil when there is none.
func FindProfile(name string) *Profile {
	for _, profile := range GetAwsProfiles() {
		if profile.Name == name {
			return &profile
		}
	}
	return nil
}

func GetAwsCredentialsFilePath() string {
	if path := os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); path != "" {
		return path
	}
	return config.DefaultSharedCredentialsFilename()
}

func GetAwsProfileFilePath() string {
	if path := os.Getenv("AWS_CONFIG_FILE"); path != "" {
		return path
	}
	return config.DefaultSharedConfigFilename()
}

func loadAwsProfiles() {
	credentialsFile := loadIniFile(GetAwsCredentialsFilePath())
	configFile := loadIniFile(GetAwsProfileFilePath())
	awsProfiles = nil
	addCredentialsProfiles(credentialsFile, configFile)
	addSSOProfiles(configFile)
}

// loadIniFile returns nil, after logging why, when the file cannot be read.
func loadIniFile(path string) *ini.File {
	file, err := ini.Load(path)
	if err != nil {
		log.Println(err)
	}
	return file
}

// addCredentialsProfiles adds every profile of the credentials file, with the
// region its section of the config file gives it.
func addCredentialsProfiles(credentialsFile, configFile *ini.File) {
	if credentialsFile == nil {
		return
	}
	for _, section := range credentialsFile.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}
		addProfile(section.Name(), configuredRegion(configFile, section.Name()))
	}
}

func configuredRegion(configFile *ini.File, profileName string) string {
	if configFile == nil {
		return ""
	}
	section, _ := configFile.GetSection(configSectionName(profileName))
	return regionOf(section)
}

func configSectionName(profileName string) string {
	if profileName == "default" {
		return profileName
	}
	return configProfilePrefix + profileName
}

// addSSOProfiles adds the profiles that only exist in the config file
// because they sign in through SSO rather than with stored credentials.
func addSSOProfiles(configFile *ini.File) {
	if configFile == nil {
		return
	}
	for _, section := range configFile.Sections() {
		profileName, isProfile := strings.CutPrefix(section.Name(), configProfilePrefix)
		if !isProfile || profileExists(profileName) || !usesSSO(section) {
			continue
		}
		addProfile(profileName, regionOf(section))
	}
}

func usesSSO(section *ini.Section) bool {
	return section.HasKey("sso_session") || section.HasKey("sso_start_url")
}

func addProfile(name, region string) {
	awsProfiles = append(awsProfiles, Profile{Name: name, Region: region})
}

func profileExists(name string) bool {
	for _, profile := range awsProfiles {
		if profile.Name == name {
			return true
		}
	}
	return false
}

func regionOf(section *ini.Section) string {
	if section == nil {
		return ""
	}
	if region, err := section.GetKey("region"); err == nil {
		return region.Value()
	}
	return ""
}
