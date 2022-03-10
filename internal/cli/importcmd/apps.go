package importcmd

import "context"

type ImportChanges struct {
	Creates int
	Updates int
	Deletes int
}

func ImportApps(cli *cli, yaml *TenantConfig, ctx context.Context) (changes ImportChanges, err error) {
	// Fetch all existing apps
	existingApps, err := cli.getAllApps(ctx)
	if err != nil {
		return ImportChanges{0, 0, 0}, err
	}

	// Map all YAML items to the corresponding Go SDK models
	newApps := mapYamlToApps(yaml)

	// Distribute in create, delete, update, and conflicts arrays
	createsMap, updatesMap, deletesMap := distributeAppOperations(existingApps, newApps)

	// Perform patch diff
	updatesMap = diffApps(updatesMap, createsMap)

	// Execute creates, updates, and deletes
	createsSlice := appsMapToSlice(createsMap)
	updatesSlice := appsMapToSlice(updatesMap)
	deletesSlice := appsMapToSlice(deletesMap)

	err = processAppOperations(cli, createsSlice, updatesSlice, deletesSlice)
	if err != nil {
		return ImportChanges{0, 0, 0}, err
	}
	return ImportChanges{len(createsSlice), len(updatesSlice), len(deletesSlice)}, nil
}

func getAllApps() {
	// Do stuff here
}

func mapYamlToApps() {
	// Do stuff here
}

func distributeAppOperations(existingApps []App, newApps []App) (creates map[string]App, updates map[string]App, deletes map[string]App) {
	creates = make(map[string]App)
	for _, newApp := range newApps {
		creates[newApp.GetName()] = newApp
	}

	updates = make(map[string]App)
	deletes = make(map[string]App)

	for _, existingApp := range existingApps {
		// If the existing app matches a new app
		if _, ok := creates[existingApp.GetName()]; ok {
			// Add it to updates
			updates[existingApp.GetName()] = existingApp
		} else {
			// Add it to the deletes
			deletes[existingApp.GetName()] = existingApp
		}
	}
	// TODO: Handle conflicts
	return creates, updates, deletes
}

func diffApps() {
	// Do stuff here
}

func appsMapToSlice() {
	// Do stuff here
}

func processAppOperations() {
	// Do stuff here
}
