function(doc) {
	if (doc.type === 'recording') {
		emit(doc._id,1);
	}
}
