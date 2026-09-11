package service

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

// SeedData creates required PocketBase collections and demo data
func SeedData(app core.App) error {
	dao := app.Dao()

	// 1. Update/Ensure 'users' collection has 'role' field
	usersCol, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	if usersCol.Schema.GetFieldByName("role") == nil {
		usersCol.Schema.AddField(&schema.SchemaField{
			Name:     "role",
			Type:     schema.FieldTypeSelect,
			Required: false,
			Options: &schema.SelectOptions{
				Values:    []string{"admin", "editor", "viewer"},
				MaxSelect: 1,
			},
		})
		if err := dao.SaveCollection(usersCol); err != nil {
			log.Printf("Failed to save users collection with role field: %v\n", err)
		}
	}

	// 2. Create 'pages' collection if not exists
	if _, err := dao.FindCollectionByNameOrId("pages"); err != nil {
		pagesCol := &models.Collection{
			Name:       "pages",
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer("@request.auth.id != ''"),
			UpdateRule: types.Pointer("@request.auth.id != ''"),
			DeleteRule: types.Pointer("@request.auth.id != ''"),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "title",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "slug",
					Type:     schema.FieldTypeText,
					Required: true,
					Unique:   true,
				},
				&schema.SchemaField{
					Name:     "content",
					Type:     schema.FieldTypeEditor,
					Required: false,
				},
				&schema.SchemaField{
					Name:     "template",
					Type:     schema.FieldTypeSelect,
					Required: true,
					Options: &schema.SelectOptions{
						Values:    []string{"default", "landing", "blog"},
						MaxSelect: 1,
					},
				},
				&schema.SchemaField{
					Name:     "status",
					Type:     schema.FieldTypeSelect,
					Required: true,
					Options: &schema.SelectOptions{
						Values:    []string{"draft", "published", "archived"},
						MaxSelect: 1,
					},
				},
				&schema.SchemaField{
					Name:     "author",
					Type:     schema.FieldTypeRelation,
					Required: false,
					Options: &schema.RelationOptions{
						CollectionId: usersCol.Id,
						MaxSelect:    types.Pointer(1),
					},
				},
				&schema.SchemaField{
					Name:     "publishedAt",
					Type:     schema.FieldTypeDate,
					Required: false,
				},
			),
		}
		if err := dao.SaveCollection(pagesCol); err != nil {
			log.Printf("Failed to create pages collection: %v\n", err)
		}
	}

	// 3. Create 'notifications' collection if not exists
	if _, err := dao.FindCollectionByNameOrId("notifications"); err != nil {
		notifCol := &models.Collection{
			Name:       "notifications",
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer("@request.auth.id != ''"),
			UpdateRule: types.Pointer("@request.auth.id != ''"),
			DeleteRule: types.Pointer("@request.auth.id != ''"),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "user",
					Type:     schema.FieldTypeRelation,
					Required: false,
					Options: &schema.RelationOptions{
						CollectionId: usersCol.Id,
						MaxSelect:    types.Pointer(1),
					},
				},
				&schema.SchemaField{
					Name:     "title",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "message",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "type",
					Type:     schema.FieldTypeSelect,
					Required: true,
					Options: &schema.SelectOptions{
						Values:    []string{"info", "warning", "error", "success"},
						MaxSelect: 1,
					},
				},
				&schema.SchemaField{
					Name:     "read",
					Type:     schema.FieldTypeBool,
					Required: false,
				},
			),
		}
		if err := dao.SaveCollection(notifCol); err != nil {
			log.Printf("Failed to create notifications collection: %v\n", err)
		}
	}

	// 4. Create 'app_settings' collection if not exists
	if _, err := dao.FindCollectionByNameOrId("app_settings"); err != nil {
		settingsCol := &models.Collection{
			Name:       "app_settings",
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer("@request.auth.id != ''"),
			UpdateRule: types.Pointer("@request.auth.id != ''"),
			DeleteRule: types.Pointer("@request.auth.id != ''"),
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "key",
					Type:     schema.FieldTypeText,
					Required: true,
					Unique:   true,
				},
				&schema.SchemaField{
					Name:     "value",
					Type:     schema.FieldTypeText,
					Required: true,
				},
			),
		}
		if err := dao.SaveCollection(settingsCol); err != nil {
			log.Printf("Failed to create app_settings collection: %v\n", err)
		}
	}

	// 5. Seed default Admin user if none exists
	if _, err := dao.FindAuthRecordByEmail("users", "admin@workfly.com"); err != nil {
		adminRecord := models.NewRecord(usersCol)
		adminRecord.SetEmail("admin@workfly.com")
		adminRecord.SetUsername("admin")
		adminRecord.Set("name", "Admin User")
		adminRecord.Set("role", "admin")
		adminRecord.SetPassword("password123")
		adminRecord.SetVerified(true)
		if err := dao.SaveRecord(adminRecord); err != nil {
			log.Printf("Failed to seed admin user: %v\n", err)
		} else {
			log.Println("Seeded initial admin user: admin@workfly.com / password123")
		}
	}

	return nil
}
