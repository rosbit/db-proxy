var name2dsn = {
	'test': 'mysql:user:password@tcp(172.16.10.240:3306)/?charset=utf8mb4',
	'163': 'mysql:user:password@tcp(192.168.0.163:3306)/?charset=utf8mb4'
}

function getDSNByName(name) {
	var dsn = name2dsn[name];
	if (dsn) {
		return dsn
	}
	return ""
}
