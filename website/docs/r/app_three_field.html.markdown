---
layout: 'okta'
page_title: 'Okta: okta_app_three_field'
sidebar_current: 'docs-okta-resource-app-three-field'
description: |-
  Creates a Three Field Application.
---

# okta_app_three_field

This resource allows you to create and configure a Three Field Application.

## Example Usage

```hcl
resource "okta_app_three_field" "example" {
  label                = "Example App"
  sign_on_url          = "https://example.com/login.html"
  sign_on_redirect_url = "https://example.com"
  reveal_password      = true
  credentials_scheme   = "EDIT_USERNAME_AND_PASSWORD"
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The display name of the Application.

- `button_selector` - (Required) Login button field CSS selector.

- `password_selector` - (Required) Login password field CSS selector.

- `username_selector` - (Required) Login username field CSS selector.

- `extra_field_selector` - (Required) Extra field CSS selector.

- `extra_field_value` - (Required) Value for extra form field.

- `url` - (Required) Login URL.

- `url_regex` - (Optional) A regex that further restricts URL to the specified regex.

- `users` - (Optional) The users assigned to the application. See `okta_app_user` for a more flexible approach.
  - `DEPRECATED`: Please replace usage with the `okta_app_user` resource.

- `groups` - (Optional) Groups associated with the application. See `okta_app_group_assignment` for a more flexible approach.
  - `DEPRECATED`: Please replace usage with the `okta_app_group_assignments` (or `okta_app_group_assignment`) resource.

- `status` - (Optional) Status of application. By default, it is `"ACTIVE"`.

- `accessibility_error_redirect_url` - (Optional) Custom error page URL.

- `accessibility_login_redirect_url` - (Optional) Custom login page for this application.

- `accessibility_self_service` - (Optional) Enable self-service. By default, it is `false`.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar.

- `hide_ios` - (Optional) Do not display application icon on mobile app.

- `hide_web` - (Optional) Do not display application icon to users.

- `user_name_template` - (Optional) Username template. Default: `"${source.login}"`

- `user_name_template_suffix` - (Optional) Username template suffix.

- `user_name_template_type` - (Optional) Username template type. Default: `"BUILT_IN"`.

- `user_name_template_push_status` - (Optional) Push username on update. Valid values: `"PUSH"` and `"DONT_PUSH"`.

- `logo` - (Optional) Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.

- `admin_note` - (Optional) Application notes for admins.

- `enduser_note` - (Optional) Application notes for end users.

- `credentials_scheme` - (Optional) Application credentials scheme. Can be set to `"EDIT_USERNAME_AND_PASSWORD"`, `"ADMIN_SETS_CREDENTIALS"`, `"EDIT_PASSWORD_ONLY"`, `"EXTERNAL_PASSWORD_SYNC"`, or `"SHARED_USERNAME_AND_PASSWORD"`.

- `reveal_password` - (Optional) Allow user to reveal password. It can not be set to `true` if `credentials_scheme` is `"ADMIN_SETS_CREDENTIALS"`, `"SHARED_USERNAME_AND_PASSWORD"` or `"EXTERNAL_PASSWORD_SYNC"`.

- `shared_username` - (Optional) Shared username, required for certain schemes.

- `shared_password` - (Optional) Shared password, required for certain schemes.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`.

- `skip_groups` - (Optional) Indicator that allows the app to skip `groups` sync (it's also can be provided during import). Default is `false`.

## Attributes Reference

- `name` - Name assigned to the application by Okta.

- `sign_on_mode` - Sign-on mode of application.

- `logo_url` - Direct link of application logo.

## Timeouts

The `timeouts` block allows you to specify custom [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions: 

- `create` - Create timeout if syncing users/groups (default 1 hour).

- `update` - Update timeout if syncing users/groups (default 1 hour).

- `read` - Read timeout if syncing users/groups (default 1 hour).

## Import

A Three Field App can be imported via the Okta ID.

```
$ terraform import okta_app_three_field.example &#60;app id&#62;
```

It's also possible to import app without groups or/and users. In this case ID may look like this:

```
$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_users/skip_groups

$ terraform import okta_app_basic_auth.example &#60;app id&#62;/skip_groups
```
