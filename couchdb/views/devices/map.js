function(doc) {
	if (doc.type === 'device') {
		emit(doc._id,1);
	}
}
