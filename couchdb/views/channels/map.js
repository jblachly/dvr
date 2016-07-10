function(doc) {
	if (doc.type === 'channel') {
		emit(doc._id,1);
	}
}
