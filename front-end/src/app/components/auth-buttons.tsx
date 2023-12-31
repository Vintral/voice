import Button from "@mui/material/Button"
import { signIn, signOut } from "auth"

export function SignIn({
  provider,
  ...props
}: { provider?: string } & React.ComponentPropsWithRef<typeof Button>) {
  return (
    <form
      action={async () => {
        "use server"
        console.log( "CALL SIGNIN" );
        //await signIn(provider)
      }}
    >
      <Button variant="outlined" {...props}>Sign In!</Button>
    </form>
  )
}

export function SignOut(props: React.ComponentPropsWithRef<typeof Button>) {
  return (
    <form
      action={async () => {
        "use server"
        await signOut()
      }}
      className="w-full"
    >
      <Button variant="outlined" className="w-full p-0" {...props}>
        Sign Out
      </Button>
    </form>
  )
}