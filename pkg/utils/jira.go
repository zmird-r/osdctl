package utils

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
)

const (
	JiraTokenConfigKey = "jira_token"
	JiraBaseURL        = "https://issues.redhat.com"
)

// GetJiraClient creates a jira client that connects to
// https://issues.redhat.com. To work, the jiraToken needs to be set in the
// config
func GetJiraClient() (*jira.Client, error) {
	if !viper.IsSet(JiraTokenConfigKey) {
		return nil, fmt.Errorf("key %s is not set in config file", JiraTokenConfigKey)
	}

	jiratoken := viper.GetString(JiraTokenConfigKey)

	tp := jira.PATAuthTransport{
		Token: jiratoken,
	}
	return jira.NewClient(tp.Client(), JiraBaseURL)
}

func GetJiraIssuesForCluster(clusterID string, externalClusterID string) ([]jira.Issue, error) {
	jiraClient, err := GetJiraClient()
	if err != nil {
		return nil, fmt.Errorf("error connecting to jira: %v", err)
	}

	jql := fmt.Sprintf(
		`(project = "OpenShift Hosted SRE Support" AND "Cluster ID" ~ "%s") 
		OR (project = "OpenShift Hosted SRE Support" AND "Cluster ID" ~ "%s") 
		ORDER BY created DESC`,
		externalClusterID,
		clusterID,
	)

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search for jira issues: %w\n", err)
	}

	return issues, nil
}

func GetJiraSupportExceptionsForOrg(organizationID string) ([]jira.Issue, error) {
	jiraClient, err := GetJiraClient()
	if err != nil {
		return nil, fmt.Errorf("error connecting to jira: %v", err)
	}

	jql := fmt.Sprintf(
		`project = "Support Exceptions" AND type = Story AND Status = Approved AND
		 Resolution = Unresolved AND "Customer Name" ~ "%s"`,
		organizationID,
	)

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search for jira issues %w", err)
	}

	return issues, nil
}
