package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAppOauth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppOauthRead,
		Schema: buildSchema(skipUsersAndGroupsSchema, map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_submit_toolbar": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Do not display application icon to users",
			},
			"grant_types": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of OAuth 2.0 grant types",
			},
			"response_types": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of OAuth 2.0 response type strings.",
			},
			"redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of URIs for use in the redirect-based flow.",
			},
			"post_logout_redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of URIs for redirection after logout",
			},
			"logo_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that references a logo for the client.",
			},
			"login_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that initiates login.",
			},
			"login_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Idp-Initiated login that the client supports, if any",
			},
			"login_scopes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of scopes to use for the request when 'login_mode' == OKTA",
			},
			"client_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI to a web page providing information about the client.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OAuth client ID",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth client secret",
			},
			"policy_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI to web page providing client policy document.",
			},
			"links": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Discoverable resources related to the app",
			},
			"groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Groups associated with the application",
				Deprecated:  "The `groups` field is now deprecated for the data source `okta_app_oauth`, please replace all uses of this with: `okta_app_group_assignments`",
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the application",
				Deprecated:  "The `users` field is now deprecated for the data source `okta_app_oauth`, please replace all uses of this with: `okta_app_user_assignments`",
			},
			"wildcard_redirect": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates if the client is allowed to use wildcard matching of redirect_uris",
			},
		}),
	}
}

func dataSourceAppOauthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("invalid OAuth app filters: %v", err)
	}
	var app *okta.OpenIdConnectApplication
	if filters.ID != "" {
		respApp, _, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, filters.ID, okta.NewOpenIdConnectApplication(), nil)
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app = respApp.(*okta.OpenIdConnectApplication)
	} else {
		re := getOktaClientFromMetadata(m).GetRequestExecutor()
		qp := &query.Params{Limit: 1, Filter: filters.Status, Q: filters.getQ()}
		req, err := re.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/apps%s", qp.String()), nil)
		if err != nil {
			return diag.Errorf("failed to list OAuth apps: %v", err)
		}
		var appList []*okta.OpenIdConnectApplication
		_, err = re.Do(ctx, req, &appList)
		if err != nil {
			return diag.Errorf("failed to list OAuth apps: %v", err)
		}
		if len(appList) < 1 {
			return diag.Errorf("no OAuth application found with provided filter: %s", filters)
		}
		if filters.Label != "" && appList[0].Label != filters.Label {
			return diag.Errorf("no OAuth application found with the provided label: %s", filters.Label)
		}
		logger(m).Info("found multiple OAuth applications with the criteria supplied, using the first one, sorted by creation date")
		app = appList[0]
	}
	err = setAppUsersIDsAndGroupsIDs(ctx, d, getOktaClientFromMetadata(m), app.Id)
	if err != nil {
		return diag.Errorf("failed to list OAuth's app groups and users: %v", err)
	}

	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	_ = d.Set("auto_submit_toolbar", app.Visibility.AutoSubmitToolbar)
	_ = d.Set("hide_ios", app.Visibility.Hide.IOS)
	_ = d.Set("hide_web", app.Visibility.Hide.Web)

	respTypes := []string{}
	grantTypes := []string{}
	redirectUris := []string{}
	postLogoutRedirectUris := []string{}

	if app.Settings.OauthClient != nil {
		_ = d.Set("type", app.Settings.OauthClient.ApplicationType)
		_ = d.Set("client_uri", app.Settings.OauthClient.ClientUri)
		_ = d.Set("logo_uri", app.Settings.OauthClient.LogoUri)
		_ = d.Set("login_uri", app.Settings.OauthClient.InitiateLoginUri)
		_ = d.Set("client_id", app.Credentials.OauthClient.ClientId)

		secret, err := getCurrentlyActiveClientSecret(ctx, m, app.Id)
		if err != nil {
			return diag.Errorf("failed to fetch OAuth client secret: %v", err)
		}
		_ = d.Set("client_secret", secret)

		_ = d.Set("policy_uri", app.Settings.OauthClient.PolicyUri)
		_ = d.Set("wildcard_redirect", app.Settings.OauthClient.WildcardRedirect)
		for i := range app.Settings.OauthClient.ResponseTypes {
			respTypes = append(respTypes, string(*app.Settings.OauthClient.ResponseTypes[i]))
		}
		for i := range app.Settings.OauthClient.GrantTypes {
			grantTypes = append(grantTypes, string(*app.Settings.OauthClient.GrantTypes[i]))
		}
		redirectUris = append(redirectUris, app.Settings.OauthClient.RedirectUris...)
		postLogoutRedirectUris = append(postLogoutRedirectUris, app.Settings.OauthClient.PostLogoutRedirectUris...)
	}

	aggMap := map[string]interface{}{
		"redirect_uris":             convertStringSliceToSet(redirectUris),
		"response_types":            convertStringSliceToSet(respTypes),
		"grant_types":               convertStringSliceToSet(grantTypes),
		"post_logout_redirect_uris": convertStringSliceToSet(postLogoutRedirectUris),
	}
	if app.Settings.OauthClient != nil &&
		app.Settings.OauthClient.IdpInitiatedLogin != nil {
		_ = d.Set("login_mode", app.Settings.OauthClient.IdpInitiatedLogin.Mode)
		aggMap["login_scopes"] = convertStringSliceToSet(app.Settings.OauthClient.IdpInitiatedLogin.DefaultScope)
	}

	err = setNonPrimitives(d, aggMap)
	if err != nil {
		return diag.Errorf("failed to set OAuth application properties: %v", err)
	}
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	return nil
}

// getCurrentlyActiveClientSecret See: https://developer.okta.com/docs/reference/api/apps/#list-client-secrets
func getCurrentlyActiveClientSecret(ctx context.Context, m interface{}, appId string) (string, error) {
	secrets, _, err := getOktaClientFromMetadata(m).Application.ListClientSecretsForApplication(ctx, appId)
	if err != nil {
		return "", err
	}

	// There can only be two client secrets. Regardless, choose the latest created active secret.
	var secretValue string
	var secret *okta.ClientSecret
	for _, s := range secrets {
		if secret == nil && s.Status == "ACTIVE" {
			secret = s
		}
		if secret != nil && s.Status == "ACTIVE" && secret.Created.Before(*s.Created) {
			secret = s
		}
	}
	if secret != nil {
		secretValue = secret.ClientSecret
	}

	return secretValue, nil
}
