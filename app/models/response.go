package models

// MetaInfo -> informasi pagination & filter
type MetaInfo struct { 
    Page   int    `json:"page"` 
    Limit  int    `json:"limit"` 
    Total  int    `json:"total"` 
    Pages  int    `json:"pages"` 
    SortBy string `json:"sortBy"` 
    Order  string `json:"order"` 
    Search string `json:"search"` 
} 

// UserResponse -> hasil akhir untuk endpoint /users 
type UserResponse struct { 
    Data []User   `json:"data"` 
    Meta MetaInfo `json:"meta"` 
} 

