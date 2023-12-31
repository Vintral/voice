import { NextAuthOptions } from "next-auth";
//import CognitoProvider from "next-auth/providers/cognito";
import GitHubProvider from "next-auth/providers/github";

export const authOptions: NextAuthOptions = {
  secret: process.env.NEXTAUTH_SECRET,
  providers: [
    // CognitoProvider({
    //   clientId: process.env.COGNITO_CLIENT_ID || "",
    //   clientSecret: process.env.COGNITO_CLIENT_SECRET || "",
    //   issuer: process.env.COGNITO_ISSUER,
    // })
    GitHubProvider({
      clientId: process.env.GITHUB_ID || "",
      clientSecret: process.env.GITHUB_SECRET || "",
    })
  ],
  pages: {
    signIn: "/auth/signin",
  },
};