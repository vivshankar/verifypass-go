package config

import "os"

var (
	/*Port             string = getOrDefault("PORT", "8080")
	Tenant           string = getOrDefault("TENANT", "https://bastion.dev.verify.ibmcloudsecurity.com")
	ClientID         string = getOrDefault("CLIENT_ID", "a2dca4db-d98d-497b-83f7-61e9297848dc")
	ClientSecret     string = getOrDefault("CLIENT_SECRET", "GKsNjEoYeA")
	RedirectURI      string = getOrDefault("REDIRECT_URI", "http://localhost:8080/oauth/callback")
	EULAPurposeID    string = getOrDefault("EULA_PURPOSE_ID", "80585ec3-997b-4080-838b-bb7e6bc1dfe2")
	EULAAccessType   string = getOrDefault("EULA_ACCESS_TYPE", "global_d123ea95-1b2b-4f96-bfdd-b8596638756g")
	ProfilePurposeID string = getOrDefault("PROFILE_PURPOSE_ID", "8ac17b0b-c9f4-47ec-af5f-9a9be65968cc")
	MFAPurposeID     string = getOrDefault("MFA_PURPOSE_ID", "8b6e2071-61d1-41b4-b61d-43d72d6a523c")
	ReadAccessType   string = getOrDefault("READ_ACCESS_TYPE", "37a77dd6-d937-4489-92c1-cbc6657d5cc3")
	NotifyAccessType string = getOrDefault("NOTIFY_ACCESS_TYPE", "79c2dba7-d60e-4e74-80f0-9d7e798e7f24")*/
	Port             string = getOrDefault("PORT", "8080")
	Tenant           string = getOrDefault("TENANT", "https://harbinger.verify.ibm.com")
	ClientID         string = getOrDefault("CLIENT_ID", "82443d7e-17e6-4999-afe3-305193cfd6b3")
	ClientSecret     string = getOrDefault("CLIENT_SECRET", "eDn1QHIqST")
	RedirectURI      string = getOrDefault("REDIRECT_URI", "http://localhost:8080/oauth/callback")
	EULAPurposeID    string = getOrDefault("EULA_PURPOSE_ID", "98b56762-398b-4116-94b5-125b5ca0d831")
	EULAAccessType   string = getOrDefault("EULA_ACCESS_TYPE", "global_d123ea95-1b2b-4f96-bfdd-b8596638756g")
	ProfilePurposeID string = getOrDefault("PROFILE_PURPOSE_ID", "da814bf7-20eb-412c-9f70-b2879b1ef9a1")
	MFAPurposeID     string = getOrDefault("MFA_PURPOSE_ID", "860ebe13-f8b5-4f16-819e-1d475eeb6db5")
	ReadAccessType   string = getOrDefault("READ_ACCESS_TYPE", "6be77a89-1dc1-488a-b315-6cfe3480e48b")
	NotifyAccessType string = getOrDefault("NOTIFY_ACCESS_TYPE", "6be77a89-1dc1-488a-b315-6cfe3480e48b")
)

func getOrDefault(name string, def string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		val = def
	}
	return val
}
