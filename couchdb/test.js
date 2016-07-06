		if (doc.disease) {
			emit(doc.disease, 1);
		}
		else if (doc.diseases) {
			for(var i=0; i<doc.diseases.length; i++) {
				emit(doc.diseases[i], 1);
			}
		}
	}
}