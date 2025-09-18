package main

import (
	"fmt"
	"path"
	"runtime/metrics"
	"strings"

	"github.com/prometheus/common/model"
)

/*
用于获取runtime.metrics中的golang prometheus指标名称以及类型
参考github.com/prometheus/client_golang@v1.18.0/prometheus/internal/go_runtime_metrics.go RuntimeMetricsToProm函数

--------- 结果 -----------
go_cgo_go_to_c_calls_calls_total  counter
go_cpu_classes_gc_mark_assist_cpu_seconds_total  counter
go_cpu_classes_gc_mark_dedicated_cpu_seconds_total  counter
go_cpu_classes_gc_mark_idle_cpu_seconds_total  counter
go_cpu_classes_gc_pause_cpu_seconds_total  counter
go_cpu_classes_gc_total_cpu_seconds_total  counter
go_cpu_classes_idle_cpu_seconds_total  counter
go_cpu_classes_scavenge_assist_cpu_seconds_total  counter
go_cpu_classes_scavenge_background_cpu_seconds_total  counter
go_cpu_classes_scavenge_total_cpu_seconds_total  counter
go_cpu_classes_total_cpu_seconds_total  counter
go_cpu_classes_user_cpu_seconds_total  counter
go_gc_cycles_automatic_gc_cycles_total  counter
go_gc_cycles_forced_gc_cycles_total  counter
go_gc_cycles_total_gc_cycles_total  counter
go_gc_gogc_percent  gauge
go_gc_gomemlimit_bytes  gauge
go_gc_heap_allocs_by_size_bytes  histogram
go_gc_heap_allocs_bytes_total  counter
go_gc_heap_allocs_objects_total  counter
go_gc_heap_frees_by_size_bytes  histogram
go_gc_heap_frees_bytes_total  counter
go_gc_heap_frees_objects_total  counter
go_gc_heap_goal_bytes  gauge
go_gc_heap_live_bytes  gauge
go_gc_heap_objects_objects  gauge
go_gc_heap_tiny_allocs_objects_total  counter
go_gc_limiter_last_enabled_gc_cycle  gauge
go_gc_pauses_seconds  histogram
go_gc_scan_globals_bytes  gauge
go_gc_scan_heap_bytes  gauge
go_gc_scan_stack_bytes  gauge
go_gc_scan_total_bytes  gauge
go_gc_stack_starting_size_bytes  gauge
go_godebug_non_default_behavior_asynctimerchan_events_total  counter
go_godebug_non_default_behavior_execerrdot_events_total  counter
go_godebug_non_default_behavior_gocachehash_events_total  counter
go_godebug_non_default_behavior_gocachetest_events_total  counter
go_godebug_non_default_behavior_gocacheverify_events_total  counter
go_godebug_non_default_behavior_gotestjsonbuildtext_events_total  counter
go_godebug_non_default_behavior_gotypesalias_events_total  counter
go_godebug_non_default_behavior_http2client_events_total  counter
go_godebug_non_default_behavior_http2server_events_total  counter
go_godebug_non_default_behavior_httplaxcontentlength_events_total  counter
go_godebug_non_default_behavior_httpmuxgo121_events_total  counter
go_godebug_non_default_behavior_httpservecontentkeepheaders_events_total  counter
go_godebug_non_default_behavior_installgoroot_events_total  counter
go_godebug_non_default_behavior_multipartmaxheaders_events_total  counter
go_godebug_non_default_behavior_multipartmaxparts_events_total  counter
go_godebug_non_default_behavior_multipathtcp_events_total  counter
go_godebug_non_default_behavior_netedns0_events_total  counter
go_godebug_non_default_behavior_panicnil_events_total  counter
go_godebug_non_default_behavior_randautoseed_events_total  counter
go_godebug_non_default_behavior_randseednop_events_total  counter
go_godebug_non_default_behavior_rsa1024min_events_total  counter
go_godebug_non_default_behavior_tarinsecurepath_events_total  counter
go_godebug_non_default_behavior_tls10server_events_total  counter
go_godebug_non_default_behavior_tls3des_events_total  counter
go_godebug_non_default_behavior_tlsmaxrsasize_events_total  counter
go_godebug_non_default_behavior_tlsrsakex_events_total  counter
go_godebug_non_default_behavior_tlsunsafeekm_events_total  counter
go_godebug_non_default_behavior_winreadlinkvolume_events_total  counter
go_godebug_non_default_behavior_winsymlink_events_total  counter
go_godebug_non_default_behavior_x509keypairleaf_events_total  counter
go_godebug_non_default_behavior_x509negativeserial_events_total  counter
go_godebug_non_default_behavior_x509rsacrt_events_total  counter
go_godebug_non_default_behavior_x509usefallbackroots_events_total  counter
go_godebug_non_default_behavior_x509usepolicies_events_total  counter
go_godebug_non_default_behavior_zipinsecurepath_events_total  counter
go_memory_classes_heap_free_bytes  gauge
go_memory_classes_heap_objects_bytes  gauge
go_memory_classes_heap_released_bytes  gauge
go_memory_classes_heap_stacks_bytes  gauge
go_memory_classes_heap_unused_bytes  gauge
go_memory_classes_metadata_mcache_free_bytes  gauge
go_memory_classes_metadata_mcache_inuse_bytes  gauge
go_memory_classes_metadata_mspan_free_bytes  gauge
go_memory_classes_metadata_mspan_inuse_bytes  gauge
go_memory_classes_metadata_other_bytes  gauge
go_memory_classes_os_stacks_bytes  gauge
go_memory_classes_other_bytes  gauge
go_memory_classes_profiling_buckets_bytes  gauge
go_memory_classes_total_bytes  gauge
go_sched_gomaxprocs_threads  gauge
go_sched_goroutines_goroutines  gauge
go_sched_latencies_seconds  histogram
go_sched_pauses_stopping_gc_seconds  histogram
go_sched_pauses_stopping_other_seconds  histogram
go_sched_pauses_total_gc_seconds  histogram
go_sched_pauses_total_other_seconds  histogram
go_sync_mutex_wait_total_seconds_total  counter
*/

func RuntimeMetricsToProm(d *metrics.Description) (string, string, string, string, bool) {
	namespace := "go"

	comp := strings.SplitN(d.Name, ":", 2)
	key := comp[0]
	unit := comp[1]

	// The last path element in the key is the name,
	// the rest is the subsystem.
	subsystem := path.Dir(key[1:] /* remove leading / */)
	name := path.Base(key)

	// subsystem is translated by replacing all / and - with _.
	subsystem = strings.ReplaceAll(subsystem, "/", "_")
	subsystem = strings.ReplaceAll(subsystem, "-", "_")

	// unit is translated assuming that the unit contains no
	// non-ASCII characters.
	unit = strings.ReplaceAll(unit, "-", "_")
	unit = strings.ReplaceAll(unit, "*", "_")
	unit = strings.ReplaceAll(unit, "/", "_per_")

	// name has - replaced with _ and is concatenated with the unit and
	// other data.
	name = strings.ReplaceAll(name, "-", "_")
	name += "_" + unit
	if d.Cumulative && d.Kind != metrics.KindFloat64Histogram {
		name += "_total"
	}

	valid := model.IsValidMetricName(model.LabelValue(namespace + "_" + subsystem + "_" + name))
	switch d.Kind {
	case metrics.KindUint64:
	case metrics.KindFloat64:
	case metrics.KindFloat64Histogram:
	default:
		valid = false
	}

	var typ string
	if d.Kind == metrics.KindFloat64Histogram {
		typ = "histogram"
	} else if d.Cumulative {
		typ = "counter"
	} else {
		typ = "gauge"
	}
	return namespace, subsystem, name, typ, valid
}

func main() {
	for _, d := range metrics.All() {
		namespace, subsystem, name, typ, valid := RuntimeMetricsToProm(&d)
		if valid {
			metricName := namespace + "_" + subsystem + "_" + name
			fmt.Printf("%s  %s\n", metricName, typ)
		}
	}
}
