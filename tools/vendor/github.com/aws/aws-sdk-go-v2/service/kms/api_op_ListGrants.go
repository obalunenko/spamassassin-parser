// Code generated by smithy-go-codegen DO NOT EDIT.

package kms

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Gets a list of all grants for the specified KMS key. You must specify the KMS
// key in all requests. You can filter the grant list by grant ID or grantee
// principal. For detailed information about grants, including grant terminology,
// see Using grants
// (https://docs.aws.amazon.com/kms/latest/developerguide/grants.html) in the Key
// Management Service Developer Guide . For examples of working with grants in
// several programming languages, see Programming grants
// (https://docs.aws.amazon.com/kms/latest/developerguide/programming-grants.html).
// The GranteePrincipal field in the ListGrants response usually contains the user
// or role designated as the grantee principal in the grant. However, when the
// grantee principal in the grant is an Amazon Web Services service, the
// GranteePrincipal field contains the service principal
// (https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_principal.html#principal-services),
// which might represent several different grantee principals. Cross-account use:
// Yes. To perform this operation on a KMS key in a different Amazon Web Services
// account, specify the key ARN in the value of the KeyId parameter. Required
// permissions: kms:ListGrants
// (https://docs.aws.amazon.com/kms/latest/developerguide/kms-api-permissions-reference.html)
// (key policy) Related operations:
//
// * CreateGrant
//
// * ListRetirableGrants
//
// *
// RetireGrant
//
// * RevokeGrant
func (c *Client) ListGrants(ctx context.Context, params *ListGrantsInput, optFns ...func(*Options)) (*ListGrantsOutput, error) {
	if params == nil {
		params = &ListGrantsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListGrants", params, optFns, c.addOperationListGrantsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListGrantsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListGrantsInput struct {

	// Returns only grants for the specified KMS key. This parameter is required.
	// Specify the key ID or key ARN of the KMS key. To specify a KMS key in a
	// different Amazon Web Services account, you must use the key ARN. For example:
	//
	// *
	// Key ID: 1234abcd-12ab-34cd-56ef-1234567890ab
	//
	// * Key ARN:
	// arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab
	//
	// To
	// get the key ID and key ARN for a KMS key, use ListKeys or DescribeKey.
	//
	// This member is required.
	KeyId *string

	// Returns only the grant with the specified grant ID. The grant ID uniquely
	// identifies the grant.
	GrantId *string

	// Returns only grants where the specified principal is the grantee principal for
	// the grant.
	GranteePrincipal *string

	// Use this parameter to specify the maximum number of items to return. When this
	// value is present, KMS does not return more than the specified number of items,
	// but it might return fewer. This value is optional. If you include a value, it
	// must be between 1 and 100, inclusive. If you do not include a value, it defaults
	// to 50.
	Limit *int32

	// Use this parameter in a subsequent request after you receive a response with
	// truncated results. Set it to the value of NextMarker from the truncated response
	// you just received.
	Marker *string

	noSmithyDocumentSerde
}

type ListGrantsOutput struct {

	// A list of grants.
	Grants []types.GrantListEntry

	// When Truncated is true, this element is present and contains the value to use
	// for the Marker parameter in a subsequent request.
	NextMarker *string

	// A flag that indicates whether there are more items in the list. When this value
	// is true, the list in this response is truncated. To get more items, pass the
	// value of the NextMarker element in thisresponse to the Marker parameter in a
	// subsequent request.
	Truncated bool

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListGrantsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpListGrants{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpListGrants{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpListGrantsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListGrants(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

// ListGrantsAPIClient is a client that implements the ListGrants operation.
type ListGrantsAPIClient interface {
	ListGrants(context.Context, *ListGrantsInput, ...func(*Options)) (*ListGrantsOutput, error)
}

var _ ListGrantsAPIClient = (*Client)(nil)

// ListGrantsPaginatorOptions is the paginator options for ListGrants
type ListGrantsPaginatorOptions struct {
	// Use this parameter to specify the maximum number of items to return. When this
	// value is present, KMS does not return more than the specified number of items,
	// but it might return fewer. This value is optional. If you include a value, it
	// must be between 1 and 100, inclusive. If you do not include a value, it defaults
	// to 50.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListGrantsPaginator is a paginator for ListGrants
type ListGrantsPaginator struct {
	options   ListGrantsPaginatorOptions
	client    ListGrantsAPIClient
	params    *ListGrantsInput
	nextToken *string
	firstPage bool
}

// NewListGrantsPaginator returns a new ListGrantsPaginator
func NewListGrantsPaginator(client ListGrantsAPIClient, params *ListGrantsInput, optFns ...func(*ListGrantsPaginatorOptions)) *ListGrantsPaginator {
	if params == nil {
		params = &ListGrantsInput{}
	}

	options := ListGrantsPaginatorOptions{}
	if params.Limit != nil {
		options.Limit = *params.Limit
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListGrantsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListGrantsPaginator) HasMorePages() bool {
	return p.firstPage || p.nextToken != nil
}

// NextPage retrieves the next ListGrants page.
func (p *ListGrantsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListGrantsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.Marker = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.Limit = limit

	result, err := p.client.ListGrants(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextMarker

	if p.options.StopOnDuplicateToken && prevToken != nil && p.nextToken != nil && *prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListGrants(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "kms",
		OperationName: "ListGrants",
	}
}
