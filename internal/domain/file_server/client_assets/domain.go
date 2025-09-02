package clientassets

import "time"

// SupportedImageTypes supported image MIME types.
var SupportedImageTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
}

type ClientAssets struct {
	ID                       int64     `json:"id"`
	Quality                  int       `json:"quality"`
	NameFile                 string    `json:"name_file"`
	NameFileWithoutExtension string    `json:"name_file_without_extension"`
	ServerPathFile           string    `json:"server_path_file"`
	ClientPathFile           string    `json:"client_path_file"`
	Extension                string    `json:"extension"`
	OldNameFile              string    `json:"old_name_file"`
	OldExtension             string    `json:"old_extension"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type UploadAndConvertToWebpResponse struct {
	Quality                  int    `json:"quality"`
	NameFile                 string `json:"name_file"`
	NameFileWithoutExtension string `json:"name_file_without_extension"`
	ServerPathFile           string `json:"server_path_file"`
	ClientPathFile           string `json:"client_path_file"`
	Extension                string `json:"extension"`
	OldNameFile              string `json:"old_name_file"`
	OldExtension             string `json:"old_extension"`
}

//
// SWAGGER
//

type CreateSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    struct {
		ID                       int64     `json:"id" example:"1"`
		Quality                  int       `json:"quality" example:"30"`
		NameFile                 string    `json:"name_file" example:"01K44X76FBXJYK4D153WHZFXH7.webp"`
		NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76FBXJYK4D153WHZFXH7"`
		ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/client/01K44X76FBXJYK4D153WHZFXH7.webp"`
		ClientPathFile           string    `json:"client_path_file" example:"/images/client/01K44X76FBXJYK4D153WHZFXH7.webp"`
		Extension                string    `json:"extension" example:".webp"`
		OldNameFile              string    `json:"old_name_file" example:"img.png"`
		OldExtension             string    `json:"old_extension" example:".png"`
		CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
	} `json:"data"`
}

type AllSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID                       int64     `json:"id" example:"1"`
		Quality                  int       `json:"quality" example:"30"`
		NameFile                 string    `json:"name_file" example:"01K44X76FBXJYK4D153WHZFXH7.webp"`
		NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76FBXJYK4D153WHZFXH7"`
		ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/client/01K44X76FBXJYK4D153WHZFXH7.webp"`
		ClientPathFile           string    `json:"client_path_file" example:"/images/client/01K44X76FBXJYK4D153WHZFXH7.webp"`
		Extension                string    `json:"extension" example:".webp"`
		OldNameFile              string    `json:"old_name_file" example:"img.png"`
		OldExtension             string    `json:"old_extension" example:".png"`
		CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
	} `json:"data"`
}

type DeleteByIDSwaggerResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message" example:"success"`
	Error   string `json:"error" example:""`
	Data    []struct {
		ID                       int64     `json:"id" example:"1"`
		Quality                  int       `json:"quality" example:"30"`
		NameFile                 string    `json:"name_file" example:"01K44X76FBXJYK4D153WHZFXH7.webp"`
		NameFileWithoutExtension string    `json:"name_file_without_extension" example:"01K44X76FBXJYK4D153WHZFXH7"`
		ServerPathFile           string    `json:"server_path_file" example:"testdata/file_server/images/client/01K44X76FBXJYK4D153WHZFXH7.webp"`
		ClientPathFile           string    `json:"client_path_file" example:"/images/client/01K44X76FBXJYK4D153WHZFXH7.webp"`
		Extension                string    `json:"extension" example:".webp"`
		OldNameFile              string    `json:"old_name_file" example:"img.png"`
		OldExtension             string    `json:"old_extension" example:".png"`
		CreatedAt                time.Time `json:"created_at" example:"2025-09-02T12:48:06.37622+03:00"`
		UpdatedAt                time.Time `json:"updated_at" example:"2025-09-02T12:48:06.37622+03:00"`
	} `json:"data"`
}

type ErrorSwaggerResponse struct {
	Status  bool        `json:"status" example:"false"`
	Message string      `json:"message" example:"some error"`
	Error   string      `json:"error" example:"some error"`
	Data    interface{} `json:"data"`
}
