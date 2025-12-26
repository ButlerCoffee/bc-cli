package templates

const AlreadyLoggedInTemplate = `You are already logged in.
Do you want to login with a different account? (Y/n): `

const LoginSuccessTemplate = `
✓ Successfully logged in!
Welcome back, {{.Username}}!
`

const AuthenticatingTemplate = `
Authenticating...`

const NotLoggedInTemplate = `You are not currently logged in.`

const LogoutSuccessTemplate = `
✓ Successfully logged out!
See you next time!
`
