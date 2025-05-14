// language/util/info_title.go

package util

func TranslateInfoTitle(countryCode string) []string {
	titles := map[string][]string{
		"US": {"BloodType", "Allergy", "Medication", "Height", "Weight", "BirthDate", "Notes"},
		"GB": {"BloodType", "Allergy", "Medication", "Height", "Weight", "BirthDate", "Notes"},
		"KR": {"혈액형", "알레르기", "복용 중인 약", "키", "체중", "생년월일", "참고사항"},
		"JP": {"血液型", "アレルギー", "服用中の薬", "身長", "体重", "生年月日", "備考"},
		"CN": {"血型", "过敏史", "正在服用的药物", "身高", "体重", "出生日期", "备注"},
		"DE": {"Blutgruppe", "Allergie", "Medikamente", "Größe", "Gewicht", "Geburtsdatum", "Notizen"},
		"FR": {"Groupe sanguin", "Allergie", "Médicaments", "Taille", "Poids", "Date de naissance", "Remarques"},
		"MX": {"Tipo de sangre", "Alergia", "Medicamentos", "Estatura", "Peso", "Fecha de nacimiento", "Notas"},
	}

	if val, ok := titles[countryCode]; ok {
		return val
	}
	
	return titles["US"]
}