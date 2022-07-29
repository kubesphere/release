package document

import (
	"regexp"
	"sort"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/release/pkg/cve"
	"k8s.io/release/pkg/notes"
)

var AreaMap = map[notes.Area]notes.Area{
	// observability
	notes.AreaAlerting:      notes.AreaObservabilityAlias,
	notes.AreaAuditing:      notes.AreaObservabilityAlias,
	notes.AreaBilling:       notes.AreaObservabilityAlias,
	notes.AreaLogging:       notes.AreaObservabilityAlias,
	notes.AreaMetering:      notes.AreaObservabilityAlias,
	notes.AreaMonitoring:    notes.AreaObservabilityAlias,
	notes.AreaNotification:  notes.AreaObservabilityAlias,
	notes.AreaObservability: notes.AreaObservabilityAlias,

	// api change
	notes.AreaAPIChange: notes.AreaAPIChangeAlias,

	// app store
	notes.AreaAppManagement: notes.AreaAppAlias,
	notes.AreaApps:          notes.AreaAppAlias,

	// devops
	notes.AreaDevOps: notes.AreaDevOpsAlias,

	// development & testing
	notes.AreaE2ETestFramework: notes.AreaDevTestAlias,
	notes.AreaInfra:            notes.AreaDevTestAlias,
	notes.AreaTesting:          notes.AreaDevTestAlias,

	// edge
	notes.AreaEdge: notes.AreaEdgeAlias,

	// Multi-tenancy & Multi-cluster

	notes.AreaIAM:          notes.AreaIAMMulticlusterAlias,
	notes.AreaMulticluster: notes.AreaIAMMulticlusterAlias,

	// Service Mesh
	notes.AreaMicroService: notes.AreaMicroServiceAlias,

	// NetWork
	notes.AreaNetWork: notes.AreaNetWorkAlias,

	// storage
	notes.AreaStorage: notes.AreaStorageAlias,

	// User Experience
	notes.AreaUI: notes.AreaUIAlias,

	// Security
	notes.AreaSecurity: notes.AreaSecurityAlias,
}

var AreaPriority = []notes.Area{
	notes.AreaAlerting,
	notes.AreaAuditing,
	notes.AreaBilling,
	notes.AreaLogging,
	notes.AreaMetering,
	notes.AreaMonitoring,
	notes.AreaNotification,
	notes.AreaObservability,
	notes.AreaAPIChange,
	notes.AreaAppManagement,
	notes.AreaApps,
	notes.AreaDevOps,
	notes.AreaE2ETestFramework,
	notes.AreaInfra,
	notes.AreaTesting,
	notes.AreaEdge,
	notes.AreaIAM,
	notes.AreaMulticluster,
	notes.AreaMicroService,
	notes.AreaNetWork,
	notes.AreaStorage,
	notes.AreaUI,
	notes.AreaSecurity,
	notes.AreaUncategorized,
}

type NoteArea struct {
	Area  notes.Area
	Notes NoteCollection
}

func mapArea(area notes.Area) notes.Area {
	newArea, ok := AreaMap[area]
	if ok {
		return newArea
	}
	return area
}

type NoteAreaCollection map[notes.Area][]NoteCategory

func NewV2(
	releaseNotes *notes.ReleaseNotes,
	previousRev, currentRev string,
) (*Document, error) {
	doc := &Document{
		NotesWithActionRequired: notes.Notes{},
		Notes:                   NoteCollection{},
		NotesV2:                 make(map[notes.Area][]NoteCategory),
		CurrentRevision:         currentRev,
		PreviousRevision:        previousRev,
	}

	stripRE := regexp.MustCompile(`^([-\*]+\s+)`)
	// processNote encapsulates the pre-processing that might happen on a note
	// text before it gets bulleted during rendering.
	processNote := func(s string) string {
		return stripRE.ReplaceAllLiteralString(s, "")
	}

	areaCategory := make(map[notes.Area]map[notes.Kind]NoteCategory)

	// kindCategory := make(map[notes.Kind]NoteCategory)
	for _, pr := range releaseNotes.History() {
		note := releaseNotes.Get(pr)

		if _, hasCVE := note.DataFields["cve"]; hasCVE {
			logrus.Infof("Release note for PR #%d has CVE vulnerability info", note.PrNumber)

			// Create a new CVE data struct for the document
			newcve := cve.CVE{}

			// Populate the struct from the raw interface
			if err := newcve.ReadRawInterface(note.DataFields["cve"]); err != nil {
				return nil, errors.Wrap(err, "reading CVE data embedded in map file")
			}

			// Verify that CVE data has the minimum fields defined
			if err := newcve.Validate(); err != nil {
				return nil, errors.Wrapf(err, "checking CVE map file for PR #%d", pr)
			}
			doc.CVEList = append(doc.CVEList, newcve)
		}

		if note.DoNotPublish {
			logrus.Debugf("skipping PR %d as (marked to not be published)", pr)
			continue
		}

		// todo not consider DuplicatedArea
		// TODO: Refactor the logic here and add testing.
		// if note.DuplicateKind {
		// 	kind := mapKind(highestPriorityKind(note.Kinds))
		// 	if existing, ok := kindCategory[kind]; ok {
		// 		*existing.NoteEntries = append(*existing.NoteEntries, processNote(note.Markdown))
		// 	} else {
		// 		kindCategory[kind] = NoteCategory{Kind: kind, NoteEntries: &notes.Notes{processNote(note.Markdown)}}
		// 	}
		// } else
		if note.ActionRequired {
			doc.NotesWithActionRequired = append(doc.NotesWithActionRequired, processNote(note.Markdown))
		} else {

			var area notes.Area
			if len(note.Areas) == 0 {
				area = "Unknown Area"
			} else {
				area = notes.Area(note.Areas[0])
			}

			var kind notes.Kind
			if len(note.Kinds) == 0 {
				kind = notes.KindUncategorized
			} else {
				kind = notes.Kind(note.Kinds[0])
			}

			// TODO use the priority first area
			kindCategory, ok := areaCategory[notes.Area(area)]
			if ok {
				category, ok := kindCategory[notes.Kind(kind)]
				if ok {
					*category.NoteEntries = append(*category.NoteEntries, processNote(note.Markdown))
				} else {
					kindCategory[notes.Kind(kind)] = NoteCategory{Kind: kind, NoteEntries: &notes.Notes{processNote(note.Markdown)}}
				}
			} else {
				areaCategory[notes.Area(area)] = make(map[notes.Kind]NoteCategory)
				areaCategory[notes.Area(area)][notes.Kind(kind)] =
					NoteCategory{Kind: kind, NoteEntries: &notes.Notes{processNote(note.Markdown)}}

			}

			// for _, kind := range note.Kinds {
			// 	mappedKind := mapKind(notes.Kind(kind))

			// 	if existing, ok := kindCategory[mappedKind]; ok {
			// 		*existing.NoteEntries = append(*existing.NoteEntries, processNote(note.Markdown))
			// 	} else {
			// 		kindCategory[mappedKind] = NoteCategory{Kind: mappedKind, NoteEntries: &notes.Notes{processNote(note.Markdown)}}
			// 	}
			// }

			// if len(note.Kinds) == 0 {
			// 	// the note has not been categorized so far
			// 	kind := notes.KindUncategorized
			// 	if existing, ok := kindCategory[kind]; ok {
			// 		*existing.NoteEntries = append(*existing.NoteEntries, processNote(note.Markdown))
			// 	} else {
			// 		kindCategory[kind] = NoteCategory{Kind: kind, NoteEntries: &notes.Notes{processNote(note.Markdown)}}
			// 	}
			// }
		}
	}

	for area, areaCate := range areaCategory {
		for _, kindCate := range areaCate {
			doc.NotesV2[notes.Area(area)] = append(doc.NotesV2[notes.Area(area)], kindCate)
		}
	}
	// Do not sort Now
	// for _, category := range kindCategory {
	// 	sort.Strings(*category.NoteEntries)
	// 	doc.Notes = append(doc.Notes, category)
	// }

	// doc.Notes.Sort(kindPriority)
	sort.Strings(doc.NotesWithActionRequired)
	return doc, nil
}
