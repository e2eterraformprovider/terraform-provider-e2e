package faas

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = sweep.TestNamePrefix

func init() {
	resource.AddTestSweepers("e2e_faas_function", &resource.Sweeper{
		Name:         "e2e_faas_function",
		F:            sweepFaasFunctions,
		Dependencies: []string{
			// Add any dependencies here if needed
		},
	})
}

func sweepFaasFunctions(region string) error {
	client, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	_ = client               // Will be used when ListFunctions is implemented
	_ = context.Background() // Will be used when ListFunctions is implemented

	// TODO: Implement ListFunctions method in goe2e.FaasService
	// Once ListFunctions is available, uncomment the following code to enable sweeping:
	//
	// functions, _, err := client.FaaS.ListFunctions(ctx, opts)
	// if err != nil {
	// 	return fmt.Errorf("error listing FaaS functions: %w", err)
	// }
	//
	// log.Printf("[DEBUG] Found %d FaaS functions in total", len(functions))
	//
	// sweptCount := 0
	// for _, fn := range functions {
	// 	if !strings.HasPrefix(fn.Name, testNamePrefix) {
	// 		log.Printf("[DEBUG] Skipping FaaS function %s (does not have test prefix)", fn.Name)
	// 		continue
	// 	}
	//
	// 	log.Printf("[INFO] Deleting FaaS function: %s (ID: %s)", fn.Name, fn.ID)
	//
	// 	_, err := client.FaaS.DeleteFunction(ctx, fn.ID, opts)
	// 	if err != nil {
	// 		log.Printf("[ERROR] Failed to delete FaaS function %s: %v", fn.Name, err)
	// 		continue
	// 	}
	//
	// 	sweptCount++
	// 	log.Printf("[INFO] Successfully deleted FaaS function: %s", fn.Name)
	// }
	//
	// if err := sweepFaasNamespaces(client, ctx, opts); err != nil {
	// 	log.Printf("[WARNING] Error sweeping namespaces: %v", err)
	// }
	//
	// log.Printf("[INFO] Swept %d FaaS functions", sweptCount)

	log.Printf("[INFO] FaaS function sweeping not yet implemented - waiting for ListFunctions API method")

	return nil
}

// sweepFaasNamespaces sweeps test FaaS namespaces
// TODO: Implement ListNamespaces method in goe2e.FaasService before enabling
func sweepFaasNamespaces(client *goe2e.Client, ctx context.Context, opts *goe2e.RequestOptions) error {
	// Once ListNamespaces is available, uncomment the following code:
	//
	// namespaces, _, err := client.FaaS.ListNamespaces(ctx, opts)
	// if err != nil {
	// 	return fmt.Errorf("error listing FaaS namespaces: %w", err)
	// }
	//
	// log.Printf("[DEBUG] Found %d FaaS namespaces in total", len(namespaces))
	//
	// sweptCount := 0
	// for _, ns := range namespaces {
	// 	if !strings.HasPrefix(ns.Name, testNamePrefix) {
	// 		log.Printf("[DEBUG] Skipping FaaS namespace %s (does not have test prefix)", ns.Name)
	// 		continue
	// 	}
	//
	// 	log.Printf("[INFO] Deleting FaaS namespace: %s", ns.Name)
	//
	// 	_, err := client.FaaS.DeleteNamespace(ctx, ns.Name, opts)
	// 	if err != nil {
	// 		log.Printf("[ERROR] Failed to delete FaaS namespace %s: %v", ns.Name, err)
	// 		continue
	// 	}
	//
	// 	sweptCount++
	// 	log.Printf("[INFO] Successfully deleted FaaS namespace: %s", ns.Name)
	// }
	//
	// log.Printf("[INFO] Swept %d FaaS namespaces", sweptCount)

	return nil
}
