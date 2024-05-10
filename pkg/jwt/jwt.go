package jwt

var UserRoleCustomer = "cpm_customer"
var UserRoleAdmin = "cpm_admin"
var UserRoleSales = "cpm_sales"

//func AuthJWT(
//	auth *logindomain.Users,
//	roles *[]logindomain.RolePermission,
//	resSales *logindomain.AccountSales,
//	resRole *logindomain.Role,
//	JWTSecretAccessToken []byte,
//) (string, error) {
//	id := uuid.New()
//	expTime := time.Now().Add(time.Minute * 60)
//	var token *jwt.Token
//
//	// Admin
//	if resRole.Name == UserRoleAdmin {
//		claims := &logindomain.JwtAdmin{
//			Sub:        auth.ID,
//			Name:       auth.Name,
//			Email:      auth.Email,
//			RoleId:     auth.RoleId,
//			Roles:      []string{resRole.Name},
//			Permission: []string{},
//			RegisteredClaims: jwt.RegisteredClaims{
//				Issuer:    "login",
//				ExpiresAt: jwt.NewNumericDate(expTime),
//				IssuedAt:  jwt.NewNumericDate(time.Now()),
//				NotBefore: jwt.NewNumericDate(time.Now()),
//				ID:        id.String(),
//			},
//		}
//
//		for _, role := range *roles {
//			claims.Permission = append(claims.Permission, role.Permission.Name)
//		}
//
//		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	} else {
//		// Sales
//		claims := &logindomain.JWTSales{
//			Sub:        auth.ID,
//			Name:       auth.Name,
//			Email:      resRole.Name,
//			RoleId:     auth.RoleId,
//			BranchId:   resSales.BranchID,
//			SalesId:    resSales.ID,
//			Roles:      []string{resRole.Name},
//			Permission: []string{},
//			RegisteredClaims: jwt.RegisteredClaims{
//				Issuer:    "login",
//				ExpiresAt: jwt.NewNumericDate(expTime),
//				IssuedAt:  jwt.NewNumericDate(time.Now()),
//				NotBefore: jwt.NewNumericDate(time.Now()),
//				ID:        id.String(),
//			},
//		}
//
//		for _, role := range *roles {
//			claims.Permission = append(claims.Permission, role.Permission.Name)
//		}
//
//		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	}
//
//	return token.SignedString(JWTSecretAccessToken)
//
//}
