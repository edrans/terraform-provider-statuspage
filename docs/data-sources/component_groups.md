---
page_title: "statuspage_component_groups Data Source - terraform-provider-statuspage"
subcategory: ""
description: |-
  
---

# Data Source `statuspage_component_groups`



## Example Usage

```terraform
data "statuspage_component_groups" "default" {
    
    page_id = local.page_id

    filter {
        name = "name"
        values = [ "value_1", "value_2" ]
    }
}
```

## Schema

### Required

- **page_id** (String, Required) the ID of the page this component belongs to

### Optional

- **filter** (Block Set) (see [below for nested schema](#nestedblock--filter))
- **id** (String, Optional) The ID of this resource.

### Read-only

- **component_groups** (List of Object, Read-only) (see [below for nested schema](#nestedatt--component_groups))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- **name** (String, Required)
- **values** (List of String, Required)

Optional:

- **regex** (Boolean, Optional)


<a id="nestedatt--component_groups"></a>
### Nested Schema for `component_groups`

- **components** (List of String)
- **description** (String)
- **id** (String)
- **name** (String)
- **position** (Number)

