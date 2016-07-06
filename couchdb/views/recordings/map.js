function(doc) {
	if (doc.type === 'recording') {
		emit(doc.id,1);
	}
}
