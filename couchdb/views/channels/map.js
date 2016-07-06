function(doc) {
	if (doc.type === 'channel') {
		emit(doc.id,1);
	}
}
