package epp

// Login authenticates and authorizes an EPP session.
// Supply a non-empty value in NewPassword to change the password for subsequent sessions.
type Login struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	ClientID    string   `xml:"command>login>clID"`
	Password    string   `xml:"command>login>pw"`
	NewPassword string   `xml:"command>login>newPW,omitempty"`
	Version     string   `xml:"command>login>options>version"`
	Language    string   `xml:"command>login>options>lang"`
	Objects     []string `xml:"command>login>svcs>objURI"`
	Extensions  []string `xml:"command>login>svcs>svcExtension>extURI,omitempty"`
	TxnID       string   `xml:"command>clTRID"`
}

// <epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
//   <command>
//     <login>
//       <clID>ClientX</clID>
//       <pw>foo-BAR2</pw>
//       <newPW>bar-FOO2</newPW>
//       <options>
//         <version>1.0</version>
//         <lang>en</lang>
//       </options>
//       <svcs>
//         <objURI>urn:ietf:params:xml:ns:obj1</objURI>
//         <objURI>urn:ietf:params:xml:ns:obj2</objURI>
//         <objURI>urn:ietf:params:xml:ns:obj3</objURI>
//         <svcExtension>
//           <extURI>http://custom/obj1ext-1.0</extURI>
//         </svcExtension>
//       </svcs>
//     </login>
//     <clTRID>ABC-12345</clTRID>
//   </command>
// </epp>

// Login initializes an authenticated EPP session.
func (c *Conn) Login(clientID, password, newPassword string) (err error) {
	msg := Login{
		ClientID:    clientID,
		Password:    password,
		NewPassword: newPassword,
		Version:     "1.0",
		Language:    "en",
		TxnID:       c.id(),
	}
	if c.Greeting != nil {
		// FIXME: find the highest protocol version?
		// Do any EPP servers send anything other than 1.0?
		if len(c.Greeting.ServiceMenu.Versions) > 0 {
			msg.Version = c.Greeting.ServiceMenu.Versions[0]
		}
		// FIXME: look for a particular language?
		// Do any EPP servers send anything other than “en”?
		if len(c.Greeting.ServiceMenu.Languages) > 0 {
			msg.Language = c.Greeting.ServiceMenu.Languages[0]
		}
		// FIXME: we currently just echo back what’s reported by the server.
		// We may or may not use any of these in a given session. Optimization opportunity?
		msg.Objects = c.Greeting.ServiceMenu.Objects
		msg.Extensions = c.Greeting.ServiceMenu.Extensions
	}
	err = c.WriteMessage(&msg)
	if err != nil {
		return
	}
	rmsg := Response{}
	return c.ReadResponse(&rmsg)
}
