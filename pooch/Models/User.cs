using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Data;
using System.Data.SqlClient;
using System.DirectoryServices;
using System.DirectoryServices.Protocols;
using System.Linq;
using System.Web;

namespace pooch.Models
{
    public class User
    {
        public string Username { get; set; }
        public string Password { get; set; }

        public bool authenticate()
        {
            bool authd = false;

            try
            {
                LdapConnection ldapConn = new LdapConnection("10.0.80.19:636");
                using (ldapConn)
                {
                    var netCred = new System.Net.NetworkCredential(this.Username, this.Password, "JAYNET"); // create network credentials using username, password, and domain
                    ldapConn.SessionOptions.SecureSocketLayer = true;                                       // we are using port 636, so this needs to be true (if using port 389 set to false)
                    ldapConn.SessionOptions.VerifyServerCertificate += delegate { return true; };           // accept the certificate (no need to check since this app runs internally on the network)
                    ldapConn.AuthType = AuthType.Basic;                                                     // non-TLS and unsecured auth (for NTLM\Kerberos use AuthType.Negotiate)
                    ldapConn.Bind(netCred);                                                                 // auth the user
                }
                authd = true;
            }
            catch (LdapException ldapEx)
            {
                authd = false;
            }

            return authd;
        }
    }
}