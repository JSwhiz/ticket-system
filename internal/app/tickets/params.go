package tickets

import "time"

type ListParams struct {
    StatusID     *int
    PriorityID   *int
    AssigneeID   *string
    DepartmentID *int
    CreatedFrom  *time.Time
    CreatedTo    *time.Time
    Search       *string
    Page         *int
    PageSize     *int
}

func (p *ListParams) ToMap() map[string]interface{} {
    m := make(map[string]interface{})

    if p.StatusID != nil {
        m["status_id"] = *p.StatusID
    }
    if p.PriorityID != nil {
        m["priority_id"] = *p.PriorityID
    }
    if p.AssigneeID != nil {
        m["assignee_id"] = *p.AssigneeID
    }
    if p.DepartmentID != nil {
        m["department_id"] = *p.DepartmentID
    }
    if p.CreatedFrom != nil {
        m["created_from"] = *p.CreatedFrom
    }
    if p.CreatedTo != nil {
        m["created_to"] = *p.CreatedTo
    }
    if p.Search != nil {
        m["search"] = *p.Search
    }
    if p.Page != nil {
        m["page"] = *p.Page
    }
    if p.PageSize != nil {
        m["page_size"] = *p.PageSize
    }
    return m
}
