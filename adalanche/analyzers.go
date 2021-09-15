package main

import (
	"fmt"
	"strings"

	"github.com/go-ini/ini"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/unicode"
)

// Interesting permissions on AD
var (
	ResetPwd                   = uuid.UUID{0x00, 0x29, 0x95, 0x70, 0x24, 0x6d, 0x11, 0xd0, 0xa7, 0x68, 0x00, 0xaa, 0x00, 0x6e, 0x05, 0x29}
	DSReplicationGetChanges    = uuid.UUID{0x11, 0x31, 0xf6, 0xaa, 0x9c, 0x07, 0x11, 0xd1, 0xf7, 0x9f, 0x00, 0xc0, 0x4f, 0xc2, 0xdc, 0xd2}
	DSReplicationGetChangesAll = uuid.UUID{0x11, 0x31, 0xf6, 0xad, 0x9c, 0x07, 0x11, 0xd1, 0xf7, 0x9f, 0x00, 0xc0, 0x4f, 0xc2, 0xdc, 0xd2}
	DSReplicationSyncronize    = uuid.UUID{0x11, 0x31, 0xf6, 0xab, 0x9c, 0x07, 0x11, 0xd1, 0xf7, 0x9f, 0x00, 0xc0, 0x4f, 0xc2, 0xdc, 0xd2}

	AttributeMember                                 = uuid.UUID{0xbf, 0x96, 0x79, 0xc0, 0x0d, 0xe6, 0x11, 0xd0, 0xa2, 0x85, 0x00, 0xaa, 0x00, 0x30, 0x49, 0xe2}
	AttributeSetGroupMembership, _                  = uuid.FromString("{BC0AC240-79A9-11D0-9020-00C04FC2D4CF}")
	AttributeSIDHistory                             = uuid.UUID{0x17, 0xeb, 0x42, 0x78, 0xd1, 0x67, 0x11, 0xd0, 0xb0, 0x02, 0x00, 0x00, 0xf8, 0x03, 0x67, 0xc1}
	AttributeAllowedToActOnBehalfOfOtherIdentity, _ = uuid.FromString("{3F78C3E5-F79A-46BD-A0B8-9D18116DDC79}")
	AttributeMSDSGroupMSAMembership                 = uuid.UUID{0x88, 0x8e, 0xed, 0xd6, 0xce, 0x04, 0xdf, 0x40, 0xb4, 0x62, 0xb8, 0xa5, 0x0e, 0x41, 0xba, 0x38}
	AttributeGPLink, _                              = uuid.FromString("{F30E3BBE-9FF0-11D1-B603-0000F80367C1}")
	AttributeMSDSKeyCredentialLink, _               = uuid.FromString("{5B47D60F-6090-40B2-9F37-2A4DE88F3063}")
	AttributeSecurityGUIDGUID, _                    = uuid.FromString("{bf967924-0de6-11d0-a285-00aa003049e2}")
	AttributeAltSecurityIdentitiesGUID, _           = uuid.FromString("{00FBF30C-91FE-11D1-AEBC-0000F80367C1}")
	AttributeProfilePathGUID, _                     = uuid.FromString("{bf967a05-0de6-11d0-a285-00aa003049e2}")
	AttributeScriptPathGUID, _                      = uuid.FromString("{bf9679a8-0de6-11d0-a285-00aa003049e2}")

	ExtendedRightCertificateEnroll, _ = uuid.FromString("0e10c968-78fb-11d2-90d4-00c04f79dc55")

	ValidateWriteSelfMembership, _ = uuid.FromString("bf9679c0-0de6-11d0-a285-00aa003049e2")
	ValidateWriteSPN, _            = uuid.FromString("f3a64788-5306-11d1-a9c5-0000f80367c1")

	ObjectGuidUser               = uuid.UUID{0xbf, 0x96, 0x7a, 0xba, 0x0d, 0xe6, 0x11, 0xd0, 0xa2, 0x85, 0x00, 0xaa, 0x00, 0x30, 0x49, 0xe2}
	ObjectGuidComputer           = uuid.UUID{0xbf, 0x96, 0x7a, 0x86, 0x0d, 0xe6, 0x11, 0xd0, 0xa2, 0x85, 0x00, 0xaa, 0x00, 0x30, 0x49, 0xe2}
	ObjectGuidGroup              = uuid.UUID{0xbf, 0x96, 0x7a, 0x9c, 0x0d, 0xe6, 0x11, 0xd0, 0xa2, 0x85, 0x00, 0xaa, 0x00, 0x30, 0x49, 0xe2}
	ObjectGuidDomain             = uuid.UUID{0x19, 0x19, 0x5a, 0x5a, 0x6d, 0xa0, 0x11, 0xd0, 0xaf, 0xd3, 0x00, 0xc0, 0x4f, 0xd9, 0x30, 0xc9}
	ObjectGuidGPO                = uuid.UUID{0xf3, 0x0e, 0x3b, 0xc2, 0x9f, 0xf0, 0x11, 0xd1, 0xb6, 0x03, 0x00, 0x00, 0xf8, 0x03, 0x67, 0xc1}
	ObjectGuidOU                 = uuid.UUID{0xbf, 0x96, 0x7a, 0xa5, 0x0d, 0xe6, 0x11, 0xd0, 0xa2, 0x85, 0x00, 0xaa, 0x00, 0x30, 0x49, 0xe2}
	ObjectGuidAttributeSchema, _ = uuid.FromString("{BF967A80-0DE6-11D0-A285-00AA003049E2}")

	NullGUID    = uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	UnknownGUID = uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	OwnerSID, _        = SIDFromString("S-1-3-4")
	SystemSID, _       = SIDFromString("S-1-5-18")
	CreatorOwnerSID, _ = SIDFromString("S-1-3-0")
	SelfSID, _         = SIDFromString("S-1-5-10")

	AccountOperatorsSID, _          = SIDFromString("S-1-5-32-548")
	DAdministratorSID, _            = SIDFromString("S-1-5-21domain-500")
	DAdministratorsSID, _           = SIDFromString("S-1-5-32-544")
	BackupOperatorsSID, _           = SIDFromString("S-1-5-32-551")
	DomainAdminsSID, _              = SIDFromString("S-1-5-21domain-512")
	DomainControllersSID, _         = SIDFromString("S-1-5-21domain-516")
	EnterpriseAdminsSID, _          = SIDFromString("S-1-5-21root domain-519")
	KrbtgtSID, _                    = SIDFromString("S-1-5-21domain-502")
	PrintOperatorsSID, _            = SIDFromString("S-1-5-32-550")
	ReadOnlyDomainControllersSID, _ = SIDFromString("S-1-5-21domain-521")
	SchemaAdminsSID, _              = SIDFromString("S-1-5-21root domain-518")
	ServerOperatorsSID, _           = SIDFromString("S-1-5-32-549")
)

var PwnAnalyzers = []PwnAnalyzer{
	/* It's a Unicorn, dang ...
	{
		Method: "NullDACL",
		ObjectAnalyzer: func(o *Object) []*Object {
			var results []*Object
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return results
			}
			if sd.Control&CONTROLFLAG_DACL_PRESENT != 0 || len(sd.DACL.Entries) == 0 {
				results = append(results, AllObjects.FindOrAddSID(acl.SID))
			}

			return results
		},
	}, */

	{
		Method: PwnComputerAffectedByGPO,
		ObjectAnalyzer: func(o *Object) {
			// Only for computers, you can't really pwn users this way
			if o.Type() != ObjectTypeComputer {
				return
			}
			// Find all perent containers with GP links
			var hasparent bool
			p := o
			for {
				gpoptions := p.OneAttrString(GPOptions)
				if gpoptions == "1" {
					// inheritance is blocked, so don't move upwards
					break
				}

				p, hasparent = AllObjects.Parent(p)
				if !hasparent {
					break
				}

				gplinks := strings.Trim(p.OneAttrString(GPLink), " ")
				if len(gplinks) == 0 {
					continue
				}
				// log.Debug().Msgf("GPlink for %v on container %v: %v", o.DN(), p.DN(), gplinks)
				if !strings.HasPrefix(gplinks, "[") || !strings.HasSuffix(gplinks, "]") {
					log.Error().Msgf("Error parsing gplink on %v: %v", o.DN(), gplinks)
					continue
				}
				links := strings.Split(gplinks[1:len(gplinks)-1], "][")
				for _, link := range links {
					linkinfo := strings.Split(link, ";")
					if len(linkinfo) != 2 {
						log.Error().Msgf("Error parsing gplink on %v: %v", o.DN(), gplinks)
						continue
					}
					linkedgpodn := linkinfo[0][7:] // strip LDAP:// prefix and link to this
					linktype := linkinfo[1]
					// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpol/08090b22-bc16-49f4-8e10-f27a8fb16d18
					if linktype == "1" || linktype == "3" {
						continue // Link is disabled
					}

					gpo, found := AllObjects.Find(linkedgpodn)
					if !found {
						log.Error().Msgf("Object linked to GPO that is not found %v: %v", o.DN(), linkedgpodn)
					} else {
						gpo.Pwns(o, PwnComputerAffectedByGPO, 100)
					}
				}
			}
		},
	},

	{
		Method: PwnGPOMachineConfigPartOfGPO,
		ObjectAnalyzer: func(o *Object) {
			if o.Type() != ObjectTypeContainer || o.OneAttrString(Name) != "Machine" {
				return
			}
			// Only for computers, you can't really pwn users this way
			p, hasparent := AllObjects.Parent(o)
			if !hasparent || p.Type() != ObjectTypeGroupPolicyContainer {
				if strings.Contains(p.DN(), "Policies") {
					log.Debug().Msgf("%v+", p)
				}
				return
			}
			p.Pwns(o, PwnGPOMachineConfigPartOfGPO, 100)
		},
	},
	{
		Method: PwnGPOUserConfigPartOfGPO,
		ObjectAnalyzer: func(o *Object) {
			if o.Type() != ObjectTypeContainer || o.OneAttrString(Name) != "User" {
				return
			}
			// Only for users, you can't really pwn users this way
			p, hasparent := AllObjects.Parent(o)
			if o.Type() != ObjectTypeContainer || !hasparent || p.Type() != ObjectTypeGroupPolicyContainer {
				return
			}
			p.Pwns(o, PwnGPOUserConfigPartOfGPO, 100)
		},
	},
	{
		Method: PwnCreateUser,
		ObjectAnalyzer: func(o *Object) {
			// Only for containers and org units
			if o.Type() != ObjectTypeContainer && o.Type() != ObjectTypeOrganizationalUnit {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CREATE_CHILD, ObjectGuidUser) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnCreateUser, 100)
				}
			}
			return
		},
	},
	{
		Method: PwnCreateGroup,
		ObjectAnalyzer: func(o *Object) {
			// Only for containers and org units
			if o.Type() != ObjectTypeContainer && o.Type() != ObjectTypeOrganizationalUnit {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CREATE_CHILD, ObjectGuidGroup) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnCreateGroup, 100)
				}
			}
		},
	},
	{
		Method: PwnCreateComputer,
		ObjectAnalyzer: func(o *Object) {
			// Only for containers and org units
			if o.Type() != ObjectTypeContainer && o.Type() != ObjectTypeOrganizationalUnit {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CREATE_CHILD, ObjectGuidComputer) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnCreateComputer, 100)
				}
			}
		},
	},
	{
		Method: PwnCreateAnyObject,
		ObjectAnalyzer: func(o *Object) {
			// Only for containers and org units
			if o.Type() != ObjectTypeContainer && o.Type() != ObjectTypeOrganizationalUnit {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CREATE_CHILD, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnCreateAnyObject, 100)
				}
			}
		},
	},
	{
		Method: PwnDeleteObject,
		ObjectAnalyzer: func(o *Object) {
			// Only for containers and org units
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DELETE, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnDeleteObject, 100)
				}
			}
		},
	},
	{
		Method: PwnDeleteChildrenTarget,
		ObjectAnalyzer: func(o *Object) {
			// If parent has DELETE CHILD, I can be deleted by some SID
			if parent, found := AllObjects.Find(o.ParentDN()); found {
				sd, err := parent.SecurityDescriptor()
				if err != nil {
					return
				}
				for index, acl := range sd.DACL.Entries {
					if sd.DACL.AllowObjectClass(index, parent, RIGHT_DS_DELETE_CHILD, o.ObjectTypeGUID()) {
						AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnDeleteChildrenTarget, 100)
					}
				}
			}
		},
	},
	{
		Method: PwnInheritsSecurity,
		ObjectAnalyzer: func(o *Object) {
			if sd, err := o.SecurityDescriptor(); err == nil && sd.Control&CONTROLFLAG_DACL_PROTECTED == 0 {
				pdn := o.ParentDN()
				if pdn == o.DN() {
					// just to make sure we dont loop eternally by being stupid somehow
					return
				}
				if parentobject, found := AllObjects.Find(pdn); found {
					parentobject.Pwns(o, PwnInheritsSecurity, 100)
				}
			}
		},
	},
	{
		Method: PwnMemberOfGroup,
		ObjectAnalyzer: func(o *Object) {
			// Only for groups
			if o.Type() != ObjectTypeGroup && o.Type() != ObjectTypeForeignSecurityPrincipal {
				return
			}
			// It's a group
			for _, member := range o.Members(false) {
				member.Pwns(o, PwnMemberOfGroup, 100)
			}
		},
	},
	{
		Method: PwnACLContainsDeny,
		ObjectAnalyzer: func(o *Object) {
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for _, acl := range sd.DACL.Entries {
				if acl.Type == ACETYPE_ACCESS_DENIED || acl.Type == ACETYPE_ACCESS_DENIED_OBJECT {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnACLContainsDeny, 0) // Not a probability of success, this is just an indicator
				}
			}
		},
	},
	{
		Method: PwnOwns,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			// https://www.alsid.com/crb_article/kerberos-delegation/
			// --- Citation bloc --- This is generally true, but an exception exists: positioning a Deny for the OWNER RIGHTS SID (S-1-3-4) in an object’s ACE removes the owner’s implicit control of this object’s DACL. ---------------------
			aclhasdeny := false
			for _, ace := range sd.DACL.Entries {
				if ace.Type == ACETYPE_ACCESS_DENIED && ace.SID == OwnerSID {
					aclhasdeny = true
				}
			}
			if !sd.Owner.IsNull() && !aclhasdeny {
				AllObjects.FindOrAddSID(sd.Owner).Pwns(o, PwnOwns, 100)
			}
		},
	},
	{
		Method: PwnGenericAll,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_GENERIC_ALL, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnGenericAll, 100)
				}
			}
		},
	},
	{
		Method: PwnWriteAll,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_GENERIC_WRITE, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteAll, 100)
				}
			}
		},
	},
	{
		Method: PwnWritePropertyAll,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWritePropertyAll, 100)
				}
			}
		},
	},
	{
		Method: PwnWriteExtendedAll,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY_EXTENDED, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteExtendedAll, 100)
				}
			}
		},
	},
	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/c79a383c-2b3f-4655-abe7-dcbb7ce0cfbe IMPORTANT
	{
		Method: PwnTakeOwnership,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_WRITE_OWNER, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnTakeOwnership, 100)
				}
			}
		},
	},
	{
		Method: PwnWriteDACL,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_WRITE_DACL, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteDACL, 100)
				}
			}
		},
	},
	{
		Method:      PwnWriteAttributeSecurityGUID,
		Description: `Allows an attacker to modify the attribute security set of an attribute, promoting it to a weaker attribute set`,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if o.Type() != ObjectTypeAttributeSchema {
				return
			}
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeSecurityGUIDGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteAttributeSecurityGUID, 25) // Experimental, I've never run into this misconfiguration
				}
			}
		},
	},
	{
		Method: PwnResetPassword,
		ObjectAnalyzer: func(o *Object) {
			// Only computers and users
			if o.Type() != ObjectTypeUser && o.Type() != ObjectTypeComputer && o.Type() != ObjectTypeManagedServiceAccount {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CONTROL_ACCESS, ResetPwd) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnResetPassword, 100)
				}
			}
		},
	},
	{
		Method: PwnHasSPN,
		ObjectAnalyzer: func(o *Object) {
			// Only computers and users
			if o.Type() != ObjectTypeUser {
				return
			}
			if len(o.Attr(ServicePrincipalName)) > 0 {
				o.SetAttr(MetaHasSPN, "1")
				AuthenticatedUsers, found := AllObjects.Find("CN=Authenticated Users,CN=WellKnown Security Principals,CN=Configuration," + AllObjects.Base)
				if !found {
					log.Error().Msgf("Could not locate Authenticated Users")
					return
				}
				AuthenticatedUsers.Pwns(o, PwnHasSPN, 50)
			}
		},
	},
	{
		Method: PwnHasSPNNoPreauth,
		ObjectAnalyzer: func(o *Object) {
			// Only computers and users
			if o.Type() != ObjectTypeUser {
				return
			}
			uac, ok := o.AttrInt(UserAccountControl)
			if !ok {
				return
			}
			if uac&0x400000 == 0 {
				return
			}
			if len(o.Attr(ServicePrincipalName)) > 0 {
				AttackerObject.Pwns(o, PwnHasSPNNoPreauth, 50)
			}
		},
	},
	{
		Method: PwnWriteSPN, // Same GUID as Validated writes, just a different permission (?)
		ObjectAnalyzer: func(o *Object) {
			// Only computers and users
			if o.Type() != ObjectTypeUser {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, ValidateWriteSPN) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteSPN, 30)
				}
			}
		},
	},
	{
		Method: PwnWriteValidatedSPN,
		ObjectAnalyzer: func(o *Object) {
			// Only computers and users
			if o.Type() != ObjectTypeUser {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY_EXTENDED, ValidateWriteSPN) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteValidatedSPN, 30)
				}
			}
		},
	},
	{
		Method:      PwnWriteAllowedToAct,
		Description: `Modify the msDS-AllowedToActOnBehalfOfOtherIdentity on a computer to enable any SPN enabled user to impersonate anyone else`,
		ObjectAnalyzer: func(o *Object) {
			// Only computers
			if o.Type() != ObjectTypeComputer {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeAllowedToActOnBehalfOfOtherIdentity) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteAllowedToAct, 100) // Success rate?
				}
			}
		},
	},
	{
		Method: PwnAddMember,
		ObjectAnalyzer: func(o *Object) {
			// Only for groups
			if o.Type() != ObjectTypeGroup {
				return
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeMember) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnAddMember, 100)
				}
			}
		},
	},
	{
		Method: PwnAddMemberGroupAttr,
		ObjectAnalyzer: func(o *Object) {
			// Only for groups
			if o.Type() != ObjectTypeGroup {
				return
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeSetGroupMembership) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnAddMemberGroupAttr, 100)
				}
			}
		},
	},
	{
		Method: PwnAddSelfMember,
		ObjectAnalyzer: func(o *Object) {
			// Only for groups
			if o.Type() != ObjectTypeGroup {
				return
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY_EXTENDED, ValidateWriteSelfMembership) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnAddMember, 100)
				}
			}
		},
	},
	{
		Method: PwnReadMSAPassword,
		ObjectAnalyzer: func(o *Object) {
			msasds := o.AttrString(MSDSGroupMSAMembership)
			for _, msasd := range msasds {
				sd, err := ParseSecurityDescriptor([]byte(msasd))
				if err == nil {
					for _, acl := range sd.DACL.Entries {
						if acl.Type == ACETYPE_ACCESS_ALLOWED {
							AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnReadMSAPassword, 100)
						}
					}
				}
			}
		},
	},
	{
		Method:      PwnWriteAltSecurityIdentities,
		Description: "Allows an attacker to define a certificate that can be used to authenticate as the user",
		ObjectAnalyzer: func(o *Object) {
			// Only for users
			if o.Type() != ObjectTypeUser {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeAltSecurityIdentitiesGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteAltSecurityIdentities, 100)
				}
			}
		},
	},
	{
		Method:      PwnWriteProfilePath,
		Description: "Allows an attacker to trigger a user auth against an attacker controlled UNC path (responder)",
		ObjectAnalyzer: func(o *Object) {
			// Only for users
			if o.Type() != ObjectTypeUser {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeProfilePathGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteProfilePath, 100)
				}
			}
		},
	},
	{
		Method:      PwnWriteScriptPath,
		Description: "Allows an attacker to trigger a user auth against an attacker controlled UNC path (responder)",
		ObjectAnalyzer: func(o *Object) {
			// Only for users
			if o.Type() != ObjectTypeUser {
				return
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeScriptPathGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteScriptPath, 100)
				}
			}
		},
	},
	{
		Method: PwnHasMSA,
		ObjectAnalyzer: func(o *Object) {
			msas := o.Attr(MSDSHostServiceAccount)
			for _, dn := range msas {
				if targetmsa, found := AllObjects.Find(dn.String()); found {
					o.Pwns(targetmsa, PwnHasMSA, 100)
				}
			}
		},
	},
	{
		Method: PwnWriteKeyCredentialLink,
		ObjectAnalyzer: func(o *Object) {
			// Only for groups
			if o.Type() != ObjectTypeUser && o.Type() != ObjectTypeComputer {
				return
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_WRITE_PROPERTY, AttributeMSDSKeyCredentialLink) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnWriteKeyCredentialLink, 100)
				}
			}
		},
	},
	{
		Method: PwnSIDHistoryEquality,
		ObjectAnalyzer: func(o *Object) {
			sids := o.Attr(SIDHistory)
			for _, sidval := range sids {
				if sid, ok := sidval.Raw().(SID); ok {
					target := AllObjects.FindOrAddSID(sid)
					o.Pwns(target, PwnSIDHistoryEquality, 100)
				}
			}
		},
	},
	{
		Method: PwnAllExtendedRights,
		ObjectAnalyzer: func(o *Object) {
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CONTROL_ACCESS, NullGUID) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnAllExtendedRights, 100)
				}
			}
		},
	},
	{
		Method: PwnLocalAdminRights,
		ObjectAnalyzer: func(o *Object) {
			var pairs []SIDpair

			groupsxml := o.OneAttrString(A("_gpofile/Machine/Preferences/Groups/Groups.XML"))
			if groupsxml != "" {
				pairs = GPOparseGroups(groupsxml)
			}

			groupsini := o.OneAttrString(A("_gpofile/Machine/Microsoft/Windows NT/SecEdit/GptTmpl.Inf"))
			if groupsini != "" {
				pairs = append(pairs, GPOparseGptTmplInf(groupsini)...)
			}

			for _, sidpair := range pairs {
				if sidpair.Group == "S-1-5-32-544" {
					membersid, err := SIDFromString(sidpair.Member)
					if err == nil {
						AllObjects.FindOrAddSID(membersid).Pwns(o, PwnLocalAdminRights, 100)
					} else {
						log.Warn().Msgf("Detected Local Admin, but could not parse SID %v", sidpair.Member)
					}
				}
			}
		},
	},
	{
		Method: PwnLocalRDPRights,
		ObjectAnalyzer: func(o *Object) {
			var pairs []SIDpair

			groupsxml := o.OneAttrString(A("_gpofile/Machine/Preferences/Groups/Groups.XML"))
			if groupsxml != "" {
				pairs = GPOparseGroups(groupsxml)
			}

			groupsini := o.OneAttrString(A("_gpofile/Machine/Microsoft/Windows NT/SecEdit/GptTmpl.Inf"))
			if groupsini != "" {
				pairs = append(pairs, GPOparseGptTmplInf(groupsini)...)
			}

			for _, sidpair := range pairs {
				if sidpair.Group == "S-1-5-32-555" {
					membersid, err := SIDFromString(sidpair.Member)
					if err == nil {
						AllObjects.FindOrAddSID(membersid).Pwns(o, PwnLocalRDPRights, 30)
					} else {
						log.Warn().Msgf("Detected Local RDP, but could not parse SID %v", sidpair.Member)
					}
				}
			}
		},
	},
	{
		Method: PwnLocalDCOMRights,
		ObjectAnalyzer: func(o *Object) {
			var pairs []SIDpair

			groupsxml := o.OneAttrString(A("_gpofile/Machine/Preferences/Groups/Groups.XML"))
			if groupsxml != "" {
				pairs = GPOparseGroups(groupsxml)
			}

			groupsini := o.OneAttrString(A("_gpofile/Machine/Microsoft/Windows NT/SecEdit/GptTmpl.Inf"))
			if groupsini != "" {
				pairs = append(pairs, GPOparseGptTmplInf(groupsini)...)
			}

			for _, sidpair := range pairs {
				if sidpair.Group == "S-1-5-32-562" {
					membersid, err := SIDFromString(sidpair.Member)
					if err == nil {
						AllObjects.FindOrAddSID(membersid).Pwns(o, PwnLocalDCOMRights, 50)
					} else {
						log.Warn().Msgf("Detected Local DCOM, but could not parse SID %v", sidpair.Member)
					}
				}
			}
		},
	},
	{
		Method: PwnCertificateEnroll,
		ObjectAnalyzer: func(o *Object) {
			if o.Type() != ObjectTypeCertificateTemplate {
				return
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CONTROL_ACCESS, ExtendedRightCertificateEnroll) {
					AllObjects.FindOrAddSID(acl.SID).Pwns(o, PwnCertificateEnroll, 100)
				}
			}
		},
	},
	{
		Method: PwnScheduledTaskOnUNCPath,
		ObjectAnalyzer: func(o *Object) {
			schtasksxml := o.OneAttrString(A("_gpofile/Machine/Preferences/ScheduledTasks/ScheduledTasks.XML"))
			if schtasksxml == "" {
				return
			}
			for _, task := range GPOparseScheduledTasks(schtasksxml) {
				log.Debug().Msgf("Scheduled task: %v ... FIXME!", task)
				// if sidpair.Group == "S-1-5-32-555" {
				// 	membersid, err := SIDFromString(sidpair.Member)
				// 	if err == nil {
				// 		results = append(results, AllObjects.FindOrAddSID(membersid))
				// 	} else {
				// 		log.Warn().Msgf("Detected Local RDP, but could not parse SID %v", sidpair.Member)
				// 	}
				// }
				// log.Debug().Msgf("%v", sidpair)
			}
		},
	},

	{
		Method: PwnMachineScript,
		ObjectAnalyzer: func(o *Object) {
			scripts := o.OneAttrString(A("_gpofile/Machine/Scripts/Scripts.ini"))
			if scripts == "" {
				return
			}

			utf8 := make([]byte, len(scripts)/2)
			_, _, err := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder().Transform(utf8, []byte(scripts), true)
			if err != nil {
				utf8 = []byte(scripts)
			}

			// ini.LineBreak = "\n"

			inifile, err := ini.LoadSources(ini.LoadOptions{
				SkipUnrecognizableLines: true,
			}, utf8)

			scriptnum := 0
			for {
				k1 := inifile.Section("Startup").Key(fmt.Sprintf("%vCmdLine", scriptnum))
				k2 := inifile.Section("Startup").Key(fmt.Sprintf("%vParameters", scriptnum))
				if k1.String() == "" {
					break
				}
				// Create new synthetic object
				sob := NewObject()
				sob.SetAttr(ObjectCategory, "Script")
				sob.DistinguishedName = fmt.Sprintf("CN=Startup Script %v from GPO %v,CN=synthetic", scriptnum, o.OneAttr(Name))
				sob.SetAttr(Name, "Machine startup script "+strings.Trim(k1.String()+" "+k2.String(), " "))
				AllObjects.Add(sob)
				sob.Pwns(o, PwnMachineScript, 100)
				scriptnum++
			}

			scriptnum = 0
			for {
				k1 := inifile.Section("Shutdown").Key(fmt.Sprintf("%vCmdLine", scriptnum))
				k2 := inifile.Section("Shutdown").Key(fmt.Sprintf("%vParameters", scriptnum))
				if k1.String() == "" {
					break
				}
				// Create new synthetic object
				sob := NewObject()
				sob.DistinguishedName = fmt.Sprintf("CN=Shutdown Script %v from GPO %v,CN=synthetic", scriptnum, o.OneAttr(Name))
				sob.SetAttr(ObjectCategory, "Script")
				sob.SetAttr(Name, "Machine shutdown script "+strings.Trim(k1.String()+" "+k2.String(), " "))
				AllObjects.Add(sob)
				sob.Pwns(o, PwnMachineScript, 100)
				scriptnum++
			}
			return
		},
	},

	// LAPS password moved to pre-processing, as the attributes have different GUIDs from AD to AD (sigh)
	{
		Method: PwnDCReplicationGetChanges,
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CONTROL_ACCESS, DSReplicationGetChanges) {
					po := AllObjects.FindOrAddSID(acl.SID)
					info := dcsyncobjects[po]
					info.changes = true
					dcsyncobjects[po] = info
				}
			}
			return
		},
	},
	{
		Method: PwnDCReplicationSyncronize, // FIXME
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CONTROL_ACCESS, DSReplicationSyncronize) {
					po := AllObjects.FindOrAddSID(acl.SID)
					info := dcsyncobjects[po]
					info.sync = true
					dcsyncobjects[po] = info
				}
			}
		},
	},
	{
		Method: PwnDSReplicationGetChangesAll, // FIXME
		ObjectAnalyzer: func(o *Object) {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.AllowObjectClass(index, o, RIGHT_DS_CONTROL_ACCESS, DSReplicationGetChangesAll) {
					po := AllObjects.FindOrAddSID(acl.SID)
					info := dcsyncobjects[po]
					info.all = true
					dcsyncobjects[po] = info
				}
			}
		},
	},
}

func MakeAdminSDHolderPwnanalyzerFunc(adminsdholder *Object, excluded string) PwnAnalyzer {
	return PwnAnalyzer{
		Method: PwnAdminSDHolderOverwriteACL,
		ObjectAnalyzer: func(o *Object) {
			return // FIXME

			// Check if object is a user account
			// if o.Type() == ObjectTypeGroup {
			// Let's see if this is a protected group
			// if o.SID() == AccountOperators
			// }

			if o.Type() != ObjectTypeUser {
				return
			}
			// Check if object is member of one of the protected groups
			// mo := o.Attr(MemberOf)

			//if ac, ok := o.AttrInt(AdminCount); ok && ac > 0 {
			// This object has an AdminCount with a value more than zero, so it kinda can be pwned by the AdminSDHolder container
			adminsdholder.Pwns(o, PwnAdminSDHolderOverwriteACL, 100)
		},
	}
}

// Objects that can DC sync
type syncinfo struct {
	changes bool
	sync    bool
	all     bool
}

var dcsyncobjects = make(map[*Object]syncinfo)

type GraphObject struct {
	*Object
	Target    bool
	CanExpand int
}

type PwnGraph struct {
	Nodes       []GraphObject
	Connections []PwnConnection // Connection to Methods map
}

type PwnPair struct {
	Source, Target *Object
}

type PwnConnection struct {
	Source, Target *Object
	PwnMethodsAndProbabilities
}

func AnalyzeObjects(includeobjects, excludeobjects *Objects, methods PwnMethodBitmap, mode string, maxdepth, maxoutgoingconnections, minprobability int) (pg PwnGraph) {
	connectionsmap := make(map[PwnPair]PwnMethodsAndProbabilities) // Pwn Connection between objects
	implicatedobjectsmap := make(map[*Object]int)                  // Object -> Processed in round n
	canexpand := make(map[*Object]int)

	// Direction to search, forward = who can pwn interestingobjects, !forward = who can interstingobjects pwn
	forward := strings.HasPrefix(mode, "normal")
	// Backlinks = include all links, don't limit per round
	backlinks := strings.HasSuffix(mode, "backlinks")

	// Convert to our working map
	for _, object := range includeobjects.AsArray() {
		implicatedobjectsmap[object] = 0
	}

	somethingprocessed := true
	processinground := 1
	for somethingprocessed && maxdepth >= processinground {
		somethingprocessed = false
		log.Debug().Msgf("Processing round %v with %v total objects", processinground, len(implicatedobjectsmap))
		newimplicatedobjects := make(map[*Object]struct{})

		for object, processed := range implicatedobjectsmap {
			if processed != 0 {
				continue
			}
			somethingprocessed = true

			newconnectionsmap := make(map[PwnPair]PwnMethodsAndProbabilities) // Pwn Connection between objects

			var pwnlist PwnConnections
			if forward {
				pwnlist = object.PwnableBy
			} else {
				pwnlist = object.CanPwn
			}

			// Iterate over ever outgoing pwn
			// This is not efficient, but we sort the pwnlist first
			for _, pwntarget := range pwnlist.Objects() {
				pwninfo := pwnlist[pwntarget]

				// If this is not a chosen method, skip it
				detectedmethods := pwninfo.Intersect(methods)
				if detectedmethods.Count() == 0 || detectedmethods.IsSet(PwnACLContainsDeny) {
					// Nothing useful or just a deny ACL, skip it
					continue
				}

				// If we allow backlinks, all pwns are mapped, no matter who is the victim
				// Targets are allowed to pwn each other as a way to reach the goal of pwning all of them
				// If pwner is already processed, we don't care what it can pwn someone more far away from targets
				// If pwner is our attacker, we always want to know what it can do
				targetprocessinground, found := implicatedobjectsmap[pwntarget]
				if pwntarget != AttackerObject &&
					!backlinks &&
					found &&
					targetprocessinground != 0 &&
					targetprocessinground < processinground {
					// skip it
					continue
				}

				if excludeobjects != nil && excludeobjects.Contains(pwntarget) {
					// skip excluded objects
					// log.Debug().Msgf("Excluding target %v", pwntarget.DN())
					continue
				}

				var filteredmethods PwnMethodsAndProbabilities
				if detectedmethods == pwninfo.PwnMethodBitmap {
					filteredmethods = pwninfo
				} else {
					for _, method := range detectedmethods.Methods() {
						filteredmethods.Set(method, pwninfo.GetProbability(method)) // Sloooow
					}

				}

				newconnectionsmap[PwnPair{Source: object, Target: pwntarget}] = filteredmethods
			}

			if maxoutgoingconnections == 0 || len(newconnectionsmap) < maxoutgoingconnections {
				for pwnpair, detectedmethods := range newconnectionsmap {
					connectionsmap[pwnpair] = detectedmethods
					if _, found := implicatedobjectsmap[pwnpair.Target]; !found {
						newimplicatedobjects[pwnpair.Target] = struct{}{} // Add this to work map as non-processed
					}
				}
				// Add pwn target to graph for processing
			} else {
				log.Debug().Msgf("Outgoing expansion limit hit %v for object %v, there was %v connections", maxoutgoingconnections, object.Label(), len(newconnectionsmap))
				var addedanyway int
				for pwnpair, detectedmethods := range newconnectionsmap {
					// We assume the number of groups are limited and add them anyway
					if pwnpair.Target.Type() == ObjectTypeGroup {
						connectionsmap[pwnpair] = detectedmethods
						if _, found := implicatedobjectsmap[pwnpair.Target]; !found {
							newimplicatedobjects[pwnpair.Target] = struct{}{} // Add this to work map as non-processed
						}
						addedanyway++
					}
				}
				canexpand[object] = len(newconnectionsmap) - addedanyway
			}
			implicatedobjectsmap[object] = processinground // We're done processing this
		}
		log.Debug().Msgf("Processing round %v yielded %v new objects", processinground, len(newimplicatedobjects))
		for newentry := range newimplicatedobjects {
			implicatedobjectsmap[newentry] = 0
		}
		processinground++
	}

	// Convert map to slice
	pg.Connections = make([]PwnConnection, len(connectionsmap))
	i := 0
	for connection, methods := range connectionsmap {
		nc := PwnConnection{
			Source:                     connection.Source,
			Target:                     connection.Target,
			PwnMethodsAndProbabilities: methods,
		}
		if forward {
			nc.Source, nc.Target = nc.Target, nc.Source // swap 'em to get arrows pointing correctly
		}
		pg.Connections[i] = nc
		i++
	}

	pg.Nodes = make([]GraphObject, len(implicatedobjectsmap))
	i = 0
	for object := range implicatedobjectsmap {
		pg.Nodes[i].Object = object
		if includeobjects.Contains(object) {
			pg.Nodes[i].Target = true
		}
		if expandnum, found := canexpand[object]; found {
			pg.Nodes[i].CanExpand = expandnum
		}
		i++
	}

	return
}