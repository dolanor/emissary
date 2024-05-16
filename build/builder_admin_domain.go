package build

import (
	"bytes"
	"html/template"
	"net/http"
	"sort"

	"github.com/EmissarySocial/emissary/model"
	"github.com/EmissarySocial/emissary/service"
	"github.com/EmissarySocial/emissary/service/providers"
	"github.com/EmissarySocial/emissary/tools/dataset"
	"github.com/benpate/data"
	"github.com/benpate/derp"
	"github.com/benpate/form"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/schema"
	"github.com/benpate/rosetta/slice"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Domain is the builder for the admin/domain page
// It can only be accessed by a Domain Owner
type Domain struct {
	_provider *service.Provider
	_domain   *model.Domain

	CommonWithTemplate
}

// NewDomain returns a fully initialized `Domain` builder.
func NewDomain(factory Factory, request *http.Request, response http.ResponseWriter, template model.Template, actionID string) (Domain, error) {

	const location = "build.NewDomain"

	// Create the common Builder
	common, err := NewCommonWithTemplate(factory, request, response, template, actionID)

	if err != nil {
		return Domain{}, derp.Wrap(err, location, "Error creating common builder")
	}

	// Verify that the user is a Domain Owner
	if !common._authorization.DomainOwner {
		return Domain{}, derp.NewForbiddenError(location, "Must be domain owner to continue")
	}

	// Create and return the Domain builder
	result := Domain{
		_provider:          factory.Provider(),
		CommonWithTemplate: common,
	}

	// Find/Create new database record for the domain.
	domainService := factory.Domain()
	if _, err := domainService.LoadDomain(); err != nil {
		return Domain{}, derp.Wrap(err, location, "Error creating a new Domain")
	}

	result._domain = domainService.GetPointer()
	return result, nil
}

/******************************************
 * Renderer Interface
 ******************************************/

// Render generates the string value for this Stream
func (w Domain) Render() (template.HTML, error) {

	var buffer bytes.Buffer

	// Execute step (write HTML to buffer, update context)
	status := Pipeline(w._action.Steps).Get(w._factory, &w, &buffer)

	if status.Error != nil {
		err := derp.Wrap(status.Error, "build.Domain.Render", "Error generating HTML")
		derp.Report(err)
		return "", err
	}

	// Success!
	status.Apply(w._response)
	return template.HTML(buffer.String()), nil
}

// View executes a separate view for this Group
func (w Domain) View(actionID string) (template.HTML, error) {

	const location = "build.Domain.View"

	builder, err := NewDomain(w._factory, w._request, w._response, w._template, actionID)

	if err != nil {
		return template.HTML(""), derp.Wrap(err, location, "Error creating Group builder")
	}

	return builder.Render()
}

func (w Domain) Token() string {
	return list.Second(w.PathList())
	// return w.context().Param("param1")
}

func (w Domain) object() data.Object {
	return w._domain
}

func (w Domain) objectID() primitive.ObjectID {
	return w._domain.DomainID
}

func (w Domain) objectType() string {
	return "Domain"
}

func (w Domain) schema() schema.Schema {
	theme := w.Theme(w.ThemeID())
	return theme.Schema
}

func (w Domain) service() service.ModelService {
	return w._factory.Domain()
}

func (w Domain) NavigationID() string {
	return "admin"
}

func (w Domain) Permalink() string {
	return w.Host() + "/admin/domains"
}

func (w Domain) BasePath() string {
	return "/admin/domains"
}

func (w Domain) PageTitle() string {
	return "Settings"
}

func (w Domain) clone(action string) (Builder, error) {
	return NewDomain(w._factory, w._request, w._response, w._template, action)
}

/******************************************
 * Other Data Accessors
 ******************************************/

func (w Domain) ThemeID() string {
	return w._domain.ThemeID
}

func (w Domain) Theme(themeID string) model.Theme {
	themeService := w._factory.Theme()
	return themeService.GetTheme(themeID)
}

// PropertyForm returns the custom property form for this Domain,
// defined by the Theme that it uses.
func (w Domain) PropertyForm() form.Element {
	return w.Theme(w.ThemeID()).Form
}

/******************************************
 * Registration Methods
 ******************************************/

// RegistrationTemplates returns all available signup templates
func (w Domain) RegistrationTemplates() []form.LookupCode {
	return w._factory.Registration().List()
}

// RegistrationTemplate returns the signup template selected for this domain
func (w Domain) RegistrationTemplate() model.Registration {

	if templateID := w.QueryParam("templateId"); templateID != "" {
		if template, err := w._factory.Registration().Load(templateID); err == nil {
			return template
		}
	}

	return model.NewRegistration("", nil)
}

/******************************************
 * Other Methods
 ******************************************/

func (w Domain) Themes() []model.Theme {
	themeService := w._factory.Theme()
	result := themeService.List()

	sort.Slice(result, func(i, j int) bool {
		return model.SortThemes(result[i], result[j])
	})

	return result
}

// Providers lists all available external services that can be connected to this domain
func (w Domain) Providers() []form.LookupCode {

	providers := w._factory.Providers()

	return slice.Filter(dataset.Providers(), func(lookupCode form.LookupCode) bool {
		if lookupCode.Group == "MANUAL" {
			return true
		}

		provider, _ := providers.Get(lookupCode.Value)
		return !provider.IsEmpty()
	})
}

// Connection loads an external service connection from the database
func (w Domain) AllConnections() mapof.Object[model.Connection] {
	return w.factory().Connection().AllAsMap()
}

func (w Domain) Provider(providerID string) providers.Provider {
	result, _ := w._provider.GetProvider(providerID)
	return result
}

func (w Domain) debug() {
	log.Debug().Interface("object", w.object()).Msg("builder_admin_domain")
}
