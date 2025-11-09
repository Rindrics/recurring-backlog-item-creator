package main

func GetIssuesToCreate(config Config, month Month) IssuesToCreate {
	issuesToCreate := IssuesToCreate{
		Issues: []IssueToCreate{},
	}

	for _, candidate := range config.Issues {
		if candidate.IsCreationMonth(month) {
			issueToCreate := NewIssueToCreate(candidate, config.Defaults)
			issuesToCreate.Issues = append(issuesToCreate.Issues, issueToCreate)
		}
	}
	return issuesToCreate
}
