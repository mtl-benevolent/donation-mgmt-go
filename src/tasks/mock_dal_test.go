package tasks_test

import (
	"context"
	"donation-mgmt/src/dal"
	"sync"
	"sync/atomic"
)

type mockDAL struct {
	mutex *sync.Mutex

	pickedTasks []dal.Task
	ackCalled   *atomic.Bool
	nackCalled  *atomic.Bool
}

func (m *mockDAL) PickTasks(ctx context.Context, params dal.PickTasksParams) ([]dal.Task, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.pickedTasks) == 0 {
		return []dal.Task{}, nil
	}

	// Calculate how many tasks to return (min of available tasks and worker slots)
	tasksToReturn := int(params.WorkerSlots)
	if tasksToReturn > len(m.pickedTasks) {
		tasksToReturn = len(m.pickedTasks)
	}

	tasks := m.pickedTasks[:tasksToReturn]

	if tasksToReturn < len(m.pickedTasks) {
		m.pickedTasks = m.pickedTasks[tasksToReturn:]
	} else {
		m.pickedTasks = []dal.Task{}
	}

	return tasks, nil
}
func (m *mockDAL) AckTasks(ctx context.Context, ids []int64) (int64, error) {
	m.ackCalled.Store(true)
	return 1, nil
}
func (m *mockDAL) NackTask(ctx context.Context, params dal.NackTaskParams) (int64, error) {
	m.nackCalled.Store(true)
	return 1, nil
}
func (m *mockDAL) CountAuthorizedOrganizations(ctx context.Context, subject string) (int64, error) {
	panic("not implemented")
}
func (m *mockDAL) CountOrganizations(ctx context.Context) (int64, error) { panic("not implemented") }
func (m *mockDAL) CreateTask(ctx context.Context, arg dal.CreateTaskParams) (dal.Task, error) {
	panic("not implemented")
}
func (m *mockDAL) GetDonationByID(ctx context.Context, arg dal.GetDonationByIDParams) ([]dal.GetDonationByIDRow, error) {
	panic("not implemented")
}
func (m *mockDAL) GetDonationBySlug(ctx context.Context, arg dal.GetDonationBySlugParams) ([]dal.GetDonationBySlugRow, error) {
	panic("not implemented")
}
func (m *mockDAL) GetOrganizationByID(ctx context.Context, organizationid int64) (dal.Organization, error) {
	panic("not implemented")
}
func (m *mockDAL) GetOrganizationBySlug(ctx context.Context, slug string) (dal.Organization, error) {
	panic("not implemented")
}
func (m *mockDAL) GetOrganizationIDBySlug(ctx context.Context, slug string) (int64, error) {
	panic("not implemented")
}
func (m *mockDAL) GetOrganizationWithSettings(ctx context.Context, arg dal.GetOrganizationWithSettingsParams) (dal.GetOrganizationWithSettingsRow, error) {
	panic("not implemented")
}
func (m *mockDAL) GetScopedRoles(ctx context.Context, arg dal.GetScopedRolesParams) ([]dal.GetScopedRolesRow, error) {
	panic("not implemented")
}
func (m *mockDAL) GrantScopedRole(ctx context.Context, arg dal.GrantScopedRoleParams) (dal.ScopedUserRole, error) {
	panic("not implemented")
}
func (m *mockDAL) HasCapabilitiesForOrgByID(ctx context.Context, arg dal.HasCapabilitiesForOrgByIDParams) (dal.HasCapabilitiesForOrgByIDRow, error) {
	panic("not implemented")
}
func (m *mockDAL) HasCapabilitiesForOrgBySlug(ctx context.Context, arg dal.HasCapabilitiesForOrgBySlugParams) (dal.HasCapabilitiesForOrgBySlugRow, error) {
	panic("not implemented")
}
func (m *mockDAL) HasGlobalCapabilities(ctx context.Context, arg dal.HasGlobalCapabilitiesParams) (dal.HasGlobalCapabilitiesRow, error) {
	panic("not implemented")
}
func (m *mockDAL) InsertDonation(ctx context.Context, arg dal.InsertDonationParams) (dal.Donation, error) {
	panic("not implemented")
}
func (m *mockDAL) InsertDonationPayment(ctx context.Context, arg dal.InsertDonationPaymentParams) (dal.DonationPayment, error) {
	panic("not implemented")
}
func (m *mockDAL) InsertOrganization(ctx context.Context, arg dal.InsertOrganizationParams) (dal.Organization, error) {
	panic("not implemented")
}
func (m *mockDAL) InsertPaymentToRecurrentDonation(ctx context.Context, arg dal.InsertPaymentToRecurrentDonationParams) (dal.InsertPaymentToRecurrentDonationRow, error) {
	panic("not implemented")
}
func (m *mockDAL) ListAuthorizedOrganizations(ctx context.Context, arg dal.ListAuthorizedOrganizationsParams) ([]dal.Organization, error) {
	panic("not implemented")
}
func (m *mockDAL) ListOrganizationFiscalYears(ctx context.Context, arg dal.ListOrganizationFiscalYearsParams) ([]int16, error) {
	panic("not implemented")
}
func (m *mockDAL) ListOrganizations(ctx context.Context, arg dal.ListOrganizationsParams) ([]dal.Organization, error) {
	panic("not implemented")
}
func (m *mockDAL) RevokeScopedRoles(ctx context.Context, arg dal.RevokeScopedRolesParams) error {
	panic("not implemented")
}
func (m *mockDAL) UpdateDonationBySlug(ctx context.Context, arg dal.UpdateDonationBySlugParams) (int64, error) {
	panic("not implemented")
}
func (m *mockDAL) UpsertOrganizationSettings(ctx context.Context, arg dal.UpsertOrganizationSettingsParams) (dal.OrganizationSetting, error) {
	panic("not implemented")
}

type dalWrapper struct {
	*mockDAL
}

func (d *dalWrapper) PickTasks(ctx context.Context, params dal.PickTasksParams) ([]dal.Task, error) {
	return d.mockDAL.PickTasks(ctx, params)
}
func (d *dalWrapper) AckTasks(ctx context.Context, ids []int64) (int64, error) {
	return d.mockDAL.AckTasks(ctx, ids)
}
func (d *dalWrapper) NackTask(ctx context.Context, params dal.NackTaskParams) (int64, error) {
	return d.mockDAL.NackTask(ctx, params)
}
func (d *dalWrapper) CountAuthorizedOrganizations(ctx context.Context, subject string) (int64, error) {
	panic("not implemented")
}
func (d *dalWrapper) CountOrganizations(ctx context.Context) (int64, error) { panic("not implemented") }
func (d *dalWrapper) CreateTask(ctx context.Context, arg dal.CreateTaskParams) (dal.Task, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetDonationByID(ctx context.Context, arg dal.GetDonationByIDParams) ([]dal.GetDonationByIDRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetDonationBySlug(ctx context.Context, arg dal.GetDonationBySlugParams) ([]dal.GetDonationBySlugRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetOrganizationByID(ctx context.Context, organizationid int64) (dal.Organization, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetOrganizationBySlug(ctx context.Context, slug string) (dal.Organization, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetOrganizationIDBySlug(ctx context.Context, slug string) (int64, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetOrganizationWithSettings(ctx context.Context, arg dal.GetOrganizationWithSettingsParams) (dal.GetOrganizationWithSettingsRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) GetScopedRoles(ctx context.Context, arg dal.GetScopedRolesParams) ([]dal.GetScopedRolesRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) GrantScopedRole(ctx context.Context, arg dal.GrantScopedRoleParams) (dal.ScopedUserRole, error) {
	panic("not implemented")
}
func (d *dalWrapper) HasCapabilitiesForOrgByID(ctx context.Context, arg dal.HasCapabilitiesForOrgByIDParams) (dal.HasCapabilitiesForOrgByIDRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) HasCapabilitiesForOrgBySlug(ctx context.Context, arg dal.HasCapabilitiesForOrgBySlugParams) (dal.HasCapabilitiesForOrgBySlugRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) HasGlobalCapabilities(ctx context.Context, arg dal.HasGlobalCapabilitiesParams) (dal.HasGlobalCapabilitiesRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) InsertDonation(ctx context.Context, arg dal.InsertDonationParams) (dal.Donation, error) {
	panic("not implemented")
}
func (d *dalWrapper) InsertDonationPayment(ctx context.Context, arg dal.InsertDonationPaymentParams) (dal.DonationPayment, error) {
	panic("not implemented")
}
func (d *dalWrapper) InsertOrganization(ctx context.Context, arg dal.InsertOrganizationParams) (dal.Organization, error) {
	panic("not implemented")
}
func (d *dalWrapper) InsertPaymentToRecurrentDonation(ctx context.Context, arg dal.InsertPaymentToRecurrentDonationParams) (dal.InsertPaymentToRecurrentDonationRow, error) {
	panic("not implemented")
}
func (d *dalWrapper) ListAuthorizedOrganizations(ctx context.Context, arg dal.ListAuthorizedOrganizationsParams) ([]dal.Organization, error) {
	panic("not implemented")
}
func (d *dalWrapper) ListOrganizationFiscalYears(ctx context.Context, arg dal.ListOrganizationFiscalYearsParams) ([]int16, error) {
	panic("not implemented")
}
func (d *dalWrapper) ListOrganizations(ctx context.Context, arg dal.ListOrganizationsParams) ([]dal.Organization, error) {
	panic("not implemented")
}
func (d *dalWrapper) RevokeScopedRoles(ctx context.Context, arg dal.RevokeScopedRolesParams) error {
	panic("not implemented")
}
func (d *dalWrapper) UpdateDonationBySlug(ctx context.Context, arg dal.UpdateDonationBySlugParams) (int64, error) {
	panic("not implemented")
}
func (d *dalWrapper) UpsertOrganizationSettings(ctx context.Context, arg dal.UpsertOrganizationSettingsParams) (dal.OrganizationSetting, error) {
	panic("not implemented")
}
