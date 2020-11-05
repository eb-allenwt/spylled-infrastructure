package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/goformation/cloudformation/cognito"
	"github.com/awslabs/goformation/v4/cloudformation"

	// Name conflict, requires an explicit declaration
	// https://golang.org/ref/spec#Import_declarations
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
)

	var UserPoolName = "BuzzSpyll-Dev"
	var FederatedIdentityName = "BuzzSpyll-Dev"



func main() {

	//  Create session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	if e := createStackResources(sess); e != nil {
		fmt.Println(e.Error())

	}

}

func createStackResources(sess client.ConfigProvider) error {

	//var dynamodbBillingMode = os.Getenv("DynamodbBillingMode")



	template := cloudformation.NewTemplate()

	// tags := []tags.Tag{
	// 	tags.Tag{
	// 		Key:   "Product",
	// 		Value: "BuzzSpyll-Dev",
	// 	},
	// 	tags.Tag{
	// 		Key:   "Environment",
	// 		Value: "Dev",
	// 	},
	// }

	// Cognitio User Pool
	ro := []cognito.UserPool_RecoveryOption{}
	r1 := cognito.UserPool_RecoveryOption{
		Name:     "verified_phone_number",
		Priority: 1,
	}
	r2 := cognito.UserPool_RecoveryOption{
		Name:     "verified_email",
		Priority: 2,
	}
	ro = append(ro, r1, r2)

	template.Resources["BuzzSpyllUserPool"] = &cognito.UserPool{

		AliasAttributes: []string{"phone_number", "email", "preferred_username"},

		AccountRecoverySetting: &cognito.UserPool_AccountRecoverySetting{
			RecoveryMechanisms: ro,
		},
		AdminCreateUserConfig: &cognito.UserPool_AdminCreateUserConfig{
			AllowAdminCreateUserOnly: false,
		},
		//AutoVerifiedAttributes: []string{"email", "phone_number"},
		AutoVerifiedAttributes: []string{"email"},
		DeviceConfiguration: &cognito.UserPool_DeviceConfiguration{
			ChallengeRequiredOnNewDevice:     true,
			DeviceOnlyRememberedOnUserPrompt: false,
		},
		// EmailConfiguration
		// EmailVerificationMessage
		// EmailVerificationSubject
		// EnabledMfas
		// LambdaConfig
		MfaConfiguration: "OFF",
		Policies: &cognito.UserPool_Policies{
			PasswordPolicy: &cognito.UserPool_PasswordPolicy{
				MinimumLength:    20,
				RequireLowercase: true,
				RequireNumbers:   true,
				//RequireSymbols:                true,
				RequireUppercase:              true,
				TemporaryPasswordValidityDays: 5,
			},
		},
		// Schema
		// SmsAuthenticationMessage
		// SmsConfiguration
		// SmsVerificationMessage
		// UserPoolAddOns
		UserPoolName: UserPoolName,
		// UserPoolTags
		// UsernameAttributes: []string{"phone_number", "email"},
		// UsernameConfiguration
		// VerificationMessageTemplate

	}

	//Cognitio Federated Identity

	template.Resources["BuzzSpyllFederatedIdentity"] = &cognito.IdentityPool{

		IdentityPoolName:               FederatedIdentityName,
		AllowUnauthenticatedIdentities: true,
		AllowClassicFlow:               false,
		//CognitoIdentityProviders:       []string cognito.IdentityPool_CognitoIdentityProvider{"graph.facebook.com"}
	}

	// template.Resources["BuzzSpyllIdentityPoolRoleAttachment"] = &cognito.IdentityPoolRoleAttachment{

	// 	IdentityPoolId:             cloudformation.Ref("BuzzSpyllFederatedIdentity"),
	// 	Roles:                      rolePolicyAuth,
	// 	AWSCloudFormationDependsOn: []string{"IAMAuthRoles", "IAMUnAuthRoles"},
	// }

	// Get JSON form of AWS CloudFormation template
	j, err := template.JSON()
	if err != nil {
		fmt.Printf("Failed to generate JSON: %s\n", err)
		return err
	}
	fmt.Printf("Template creation for test is done.\n")
	//fmt.Printf("Template creation for %s Done.\n", id)
	fmt.Println("=====")
	fmt.Println("Generated template:")
	fmt.Printf("%s\n", string(j))
	fmt.Println("=====")

	// Creating the stack
	err = createStackFromBody(sess, j)
	if err != nil {
		return err
	}

	return nil
}

// Creates the stack based on values provided in createStackFromBody
func createStackFromBody(sess client.ConfigProvider, templateBody []byte) error {

	svc := cfn.New(sess)
	input := &cfn.CreateStackInput{
		TemplateBody: aws.String(string(templateBody)),
		StackName:    aws.String("BuzzSpyll-Will-1"),
		Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM")}, // Required because of creating a stack that is creating IAM resources
	}

	fmt.Println("Stack creation initiated...")

	_, err := svc.CreateStack(input)
	if err != nil {
		fmt.Println("Got error creating stack:")
		fmt.Println(err.Error())
		return err
	}

	return nil
}
