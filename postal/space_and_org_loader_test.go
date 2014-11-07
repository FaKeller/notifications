package postal_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/postal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("SpaceAndOrgLoader", func() {
    Describe("Load", func() {
        var loader postal.SpaceAndOrgLoader
        var token string
        var cc *fakes.CloudController

        BeforeEach(func() {
            cc = fakes.NewCloudController()
            cc.Spaces = map[string]cf.CloudControllerSpace{
                "space-001": cf.CloudControllerSpace{
                    GUID:             "space-001",
                    Name:             "space-name",
                    OrganizationGUID: "org-001",
                },
            }
            cc.Orgs = map[string]cf.CloudControllerOrganization{
                "org-001": cf.CloudControllerOrganization{
                    GUID: "org-001",
                    Name: "org-name",
                }, "org-123": cf.CloudControllerOrganization{
                    GUID: "org-123",
                    Name: "org-piggies",
                },
            }
            loader = postal.NewSpaceAndOrgLoader(cc)
        })

        Context("when GUID represents a space", func() {
            It("returns the name of the space and org", func() {
                space, org, err := loader.Load(postal.SpaceGUID("space-001"), token)
                if err != nil {
                    panic(err)
                }

                Expect(space).To(Equal(cf.CloudControllerSpace{
                    GUID:             "space-001",
                    Name:             "space-name",
                    OrganizationGUID: "org-001",
                }))
                Expect(org).To(Equal(cf.CloudControllerOrganization{
                    GUID: "org-001",
                    Name: "org-name",
                }))
            })

            Context("when the space cannot be found", func() {
                It("returns an error object", func() {
                    _, _, err := loader.Load(postal.SpaceGUID("space-doesnotexist"), token)

                    Expect(err).To(BeAssignableToTypeOf(postal.CCNotFoundError("")))
                    Expect(err.Error()).To(Equal(`CloudController Error: CloudController Failure (404): {"code":40004,"description":"The app space could not be found: space-doesnotexist","error_code":"CF-SpaceNotFound"}`))
                })
            })

            Context("when the org cannot be found", func() {
                It("returns an error object", func() {
                    delete(cc.Orgs, "org-001")
                    _, _, err := loader.Load(postal.SpaceGUID("space-001"), token)

                    Expect(err).To(BeAssignableToTypeOf(postal.CCNotFoundError("")))
                    Expect(err.Error()).To(Equal(`CloudController Error: CloudController Failure (404): {"code":30003,"description":"The organization could not be found: org-001","error_code":"CF-OrganizationNotFound"}`))
                })
            })

            Context("when Load returns any other type of error", func() {
                It("returns a CCDownError when the error is cf.Failure", func() {
                    failure := cf.NewFailure(401, "BOOM!")
                    cc.LoadSpaceError = failure
                    _, _, err := loader.Load(postal.SpaceGUID("space-001"), token)

                    Expect(err).To(Equal(postal.CCDownError(failure.Error())))
                })

                It("returns the same error for all other cases", func() {
                    cc.LoadSpaceError = errors.New("BOOM!")
                    _, _, err := loader.Load(postal.SpaceGUID("space-001"), token)

                    Expect(err).To(Equal(errors.New("BOOM!")))
                })
            })

            Context("when Load returns any other type of error", func() {
                It("returns a CCDownError", func() {
                    failure := cf.NewFailure(401, "BOOM!")
                    cc.LoadOrganizationError = failure
                    _, _, err := loader.Load(postal.SpaceGUID("space-001"), token)

                    Expect(err).To(Equal(postal.CCDownError(failure.Error())))
                })

                It("returns the same error for all other cases", func() {
                    cc.LoadOrganizationError = errors.New("BOOM!")
                    _, _, err := loader.Load(postal.SpaceGUID("space-001"), token)

                    Expect(err).To(Equal(errors.New("BOOM!")))
                })
            })
        })

        Context("when GUID represents an org", func() {
            It("returns the name of the org and an empty value for space", func() {
                space, org, err := loader.Load(postal.OrganizationGUID("org-123"), token)
                if err != nil {
                    panic(err)
                }

                Expect(space).To(Equal(cf.CloudControllerSpace{}))
                Expect(org).To(Equal(cf.CloudControllerOrganization{
                    GUID: "org-123",
                    Name: "org-piggies",
                }))
            })
        })

        Context("when GUID represents a user", func() {
            It("returns empty values for space, org, and error", func() {
                space, org, err := loader.Load(postal.UserGUID("user-001"), token)
                if err != nil {
                    panic(err)
                }

                Expect(space).To(Equal(cf.CloudControllerSpace{}))
                Expect(org).To(Equal(cf.CloudControllerOrganization{}))
                Expect(err).To(BeNil())
            })
        })
    })
})
