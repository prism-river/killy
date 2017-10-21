package daemon

func getAllTidb(d *Daemon) {
	data, err := d.Client.Query("count(tidb_domain_load_schema_total) by (exported_instance)")
}

func getAllTikv(d *Daemon) {
	data, err := d.Client.Query("count(tikv_engine_num_subcompaction_scheduled) by (exported_instance)")

}

func getAllPd(d *Daemon) {
	data, err := d.Client.Query("count(tikv_engine_num_subcompaction_scheduled) by (exported_instance)")
}
