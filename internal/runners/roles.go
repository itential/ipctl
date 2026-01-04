// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type RoleRunner struct {
	service *services.RoleService
	BaseRunner
	client client.Client
}

func NewRoleRunner(client client.Client, cfg *config.Config) *RoleRunner {
	return &RoleRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		client:     client,
		service:    services.NewRoleService(client),
	}
}

// GetRoles implements the `get roles` command
func (r *RoleRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	var options flags.RoleGetOptions
	utils.LoadObject(in.Options, &options)

	roles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []services.Role

	for _, ele := range roles {
		if options.All || ele.Provenance == "Custom" {
			filtered = append(filtered, ele)
		} else if options.Type != "" && ele.Provenance == options.Type {
			filtered = append(filtered, ele)
		}
	}

	return &Response{
		Keys:   []string{"name", "provenance"},
		Object: filtered,
	}, nil

}

// Describe implements the `describe role ...` command
func (r *RoleRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	var options *flags.RoleDescribeOptions
	utils.LoadObject(in.Options, &options)

	name := in.Args[0]

	var roleType string = "Custom"
	if options.Type != "" {
		roleType = options.Type
	}

	roles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var role *services.Role

	for _, ele := range roles {
		if ele.Name == name && ele.Provenance == roleType {
			role = &ele
			break
		}
	}

	if role == nil {
		return nil, errors.New(fmt.Sprintf("role `%s` does not exist", in.Args[0]))
	}

	var output = []string{
		fmt.Sprintf("Name: %s (%s), Type: %s", role.Name, role.Id, role.Provenance),
		fmt.Sprintf("Description: %s", role.Description),
	}

	return &Response{
		Text:   strings.Join(output, "\n"),
		Object: role,
	}, nil
}

// Delete implements the `delete role <name>` command
func (r *RoleRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	role, err := GetByName(r.service, in.Args[0])
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(role.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted role `%s`", in.Args[0]),
	}, nil
}

// Clear implements the `clear roles` command
func (r *RoleRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	roles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var cnt int

	for _, ele := range roles {
		if ele.Provenance == "Custom" {
			if err := r.service.Delete(ele.Id); err != nil {
				return nil, err
			}
			cnt++
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v role(s)", cnt),
	}, nil
}

func (r *RoleRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options *flags.RoleCreateOptions
	utils.LoadObject(in.Options, &options)

	role := services.Role{
		Name:       name,
		Provenance: "Custom",
	}

	availMethods, err := services.NewMethodService(r.client).GetAll()
	if err != nil {
		return nil, err
	}

	var methods []services.RoleMethod

	for _, ele := range options.AllowedMethods {
		meth, err := parseAllowedMethod(ele)
		if err != nil {
			return nil, err
		}
		if !isValidMethod(meth, availMethods) {
			return nil, errors.New(fmt.Sprintf("invalid method: name=%s, type=%s", meth.Name, meth.Provenance))
		}
		methods = append(methods, meth)
	}

	role.AllowedMethods = methods

	for _, ele := range options.AllowedViews {
		var views []services.RoleView

		val, err := parseRoleView(ele)
		if err != nil {
			return nil, err
		}
		views = append(views, val)

		role.AllowedViews = views
	}

	res, err := r.service.Create(role)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully created new role `%s`", res.Name),
	}, nil
}

func (r *RoleRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var common *flags.AssetCopyCommon
	utils.LoadObject(in.Common, &common)

	if common.From == common.To {
		return nil, errors.New("source and destination servers must be different values")
	}

	fromClient, cancel, err := NewClient(common.From, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	fromService := services.NewRoleService(fromClient)

	role, err := GetByName(fromService, name)
	if err != nil {
		return nil, err
	}

	toClient, cancel, err := NewClient(common.To, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	toService := services.NewRoleService(toClient)

	if Exists(toService, name) {
		return nil, errors.New(fmt.Sprintf("role `%s` already exists on the target server", name))
	}

	if err = isRoleValid(toClient, *role); err != nil {
		return nil, err
	}

	_, err = toService.Create(*role)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied role `%s` from `%s` to `%s`", name, common.From, common.To),
	}, nil
}

// Import implements the `import role <path>` command
func (r *RoleRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common *flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	var role services.Role
	if err := utils.ReadObjectFromDisk(path, &role); err != nil {
		return nil, err
	}

	if err := isRoleValid(r.client, role); err != nil {
		return nil, err
	}

	if common.Replace {
		existing, err := GetByName(r.service, role.Name)
		if err != nil {
			if !strings.HasSuffix(err.Error(), "does not exist") {
				return nil, err
			}
		}
		if existing != nil {
			if err := r.service.Delete(existing.Id); err != nil {
				return nil, err
			}
		}
	}

	res, err := r.service.Import(role)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported role `%s`", res.Name),
	}, nil

}

// Export implements the `export role <name>` command
func (r *RoleRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]
	roleType := "Custom"

	var common flags.AssetExportCommon
	utils.LoadObject(in.Common, &common)

	var options flags.RoleExportOptions
	utils.LoadObject(in.Options, &options)

	if options.Type != "" {
		roleType = options.Type
	}

	roles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var roleId string

	for _, ele := range roles {
		if ele.Name == name && ele.Provenance == roleType {
			roleId = ele.Id
			break
		}
	}

	if roleId == "" {
		return nil, errors.New(fmt.Sprintf("role `%s` does not exist", name))
	}

	role, err := r.service.Get(roleId)
	if err != nil {
		return nil, err
	}

	// NOTE (privatep) Deleting the ID value as it is not applicable to the
	// export and the presence of the field will cause an error when attempting
	// to import the document.  This is a byproduct of giving the illusion of
	// supporting the import/export functions in IAP
	role.Id = ""

	fn := fmt.Sprintf("%s.role.json", name)

	if err := utils.WriteJsonToDisk(role, fn, common.Path); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported role `%s` to `%s`", role.Name, fn),
	}, nil
}

func GetByName(svc *services.RoleService, name string) (*services.Role, error) {
	logger.Trace()

	roles, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range roles {
		if ele.Name == name {
			return &ele, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("role `%s` does not exist", name))
}

func Exists(svc *services.RoleService, name string) bool {
	logger.Trace()

	_, err := GetByName(svc, name)
	if err != nil {
		return false
	}

	return true
}

func parseAllowedMethod(in string) (services.RoleMethod, error) {
	elements := strings.Split(in, ",")

	var roleMethod = services.RoleMethod{}

	if len(elements) != 2 {
		return services.RoleMethod{}, errors.New(fmt.Sprintf("error parsing method: %s", in))
	}

	for _, ele := range elements {
		tokens := strings.Split(ele, "=")

		if len(tokens) != 2 {
			return services.RoleMethod{}, errors.New(fmt.Sprintf("error parsing method: %s", ele))
		}

		switch tokens[0] {
		case "name":
			roleMethod.Name = tokens[1]
		case "type":
			roleMethod.Provenance = tokens[1]
		default:
			return services.RoleMethod{}, errors.New(fmt.Sprintf("unknown keyword: %s", tokens[0]))
		}
	}

	return roleMethod, nil

}

func parseRoleView(role string) (services.RoleView, error) {
	parts := strings.Split(role, "/")
	if len(parts) != 2 {
		return services.RoleView{}, errors.New(fmt.Sprintf("error parsing ivew: %s", role))
	}
	return services.RoleView{Provenance: parts[0], Path: parts[1]}, nil
}

func isValidMethod(meth services.RoleMethod, methods []services.Method) bool {
	for _, ele := range methods {
		if ele.Name == meth.Name && ele.Provenance == meth.Provenance {
			return true
		}
	}
	return false
}

func isRoleValid(c client.Client, role services.Role) error {
	availMethods, err := services.NewMethodService(c).GetAll()
	if err != nil {
		return err
	}

	for _, ele := range role.AllowedMethods {
		if !isValidMethod(ele, availMethods) {
			return errors.New(fmt.Sprintf("configured method does not exist: name=%s, type=%s", ele.Name, ele.Provenance))
		}
	}

	return nil
}
