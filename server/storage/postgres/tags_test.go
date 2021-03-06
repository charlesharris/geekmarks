// Copyright 2017 Dmitry Frank <mail@dmitryfrank.com>
// Licensed under the BSD, see LICENSE file for details.

// +build all_tests integration_tests

package postgres

import (
	"database/sql"
	"reflect"
	"testing"

	"dmitryfrank.com/geekmarks/server/cptr"
	"dmitryfrank.com/geekmarks/server/storage"
	"dmitryfrank.com/geekmarks/server/testutils"

	"github.com/juju/errors"
)

type tagIDs struct {
	rootTagID, tag1ID, tag2ID, tag3ID, tag4ID, tag5ID, tag6ID, tag7ID, tag8ID int
}

// makeTagsHierarchy creates the following tag hierarchy for the given user:
// /
// ├── tag1
// │   └── tag3
// │       ├── tag4
// │       └── tag5
// │           └── tag6
// ├── tag2
// └── tag7
//     └── tag8
func makeTagsHierarchy(tx *sql.Tx, si *StoragePostgres, ownerID int) (ids *tagIDs, err error) {
	rootTagID, err := si.GetRootTagID(tx, ownerID)
	if err != nil {
		return nil, errors.Annotatef(err, "getting root tag for user %d", ownerID)
	}

	u1Tag1ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(rootTagID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag1", "tag1_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag1 for user %d", ownerID)
	}

	u1Tag2ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(rootTagID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag2", "tag2_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag2 for user %d", ownerID)
	}

	u1Tag3ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(u1Tag1ID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag3", "tag3_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag3 for user %d", ownerID)
	}

	u1Tag4ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(u1Tag3ID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag4_alias", "tag4"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag4 for user %d", ownerID)
	}

	u1Tag5ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(u1Tag3ID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag5", "tag5_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag5 for user %d", ownerID)
	}

	u1Tag6ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(u1Tag5ID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag6", "tag6_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag6 for user %d", ownerID)
	}

	u1Tag7ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(rootTagID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag7", "tag7_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag7 for user %d", ownerID)
	}

	u1Tag8ID, err := si.CreateTag(tx, &storage.TagData{
		OwnerID:     ownerID,
		ParentTagID: cptr.Int(u1Tag7ID),
		Description: cptr.String("test tag"),
		Names:       []string{"tag8", "tag8_alias"},
	})
	if err != nil {
		return nil, errors.Annotatef(err, "creating tag8 for user %d", ownerID)
	}

	return &tagIDs{
		rootTagID: rootTagID,
		tag1ID:    u1Tag1ID,
		tag2ID:    u1Tag2ID,
		tag3ID:    u1Tag3ID,
		tag4ID:    u1Tag4ID,
		tag5ID:    u1Tag5ID,
		tag6ID:    u1Tag6ID,
		tag7ID:    u1Tag7ID,
		tag8ID:    u1Tag8ID,
	}, nil
}

// Data created by makeTagsHierarchy
var tagsDataCreated = []storage.TagData{
	{
		ID:          2,
		OwnerID:     1,
		ParentTagID: cptr.Int(1),
		Description: cptr.String("test tag"),
		Names:       []string{"tag1", "tag1_alias"},
		Subtags: []storage.TagData{
			{
				ID:          4,
				OwnerID:     1,
				ParentTagID: cptr.Int(2),
				Description: cptr.String("test tag"),
				Names:       []string{"tag3", "tag3_alias"},
				Subtags: []storage.TagData{
					{
						ID:          5,
						OwnerID:     1,
						ParentTagID: cptr.Int(4),
						Description: cptr.String("test tag"),
						Names:       []string{"tag4_alias", "tag4"},
					},
					{
						ID:          6,
						OwnerID:     1,
						ParentTagID: cptr.Int(4),
						Description: cptr.String("test tag"),
						Names:       []string{"tag5", "tag5_alias"},
						Subtags: []storage.TagData{
							{
								ID:          7,
								OwnerID:     1,
								ParentTagID: cptr.Int(6),
								Description: cptr.String("test tag"),
								Names:       []string{"tag6", "tag6_alias"},
							},
						},
					},
				},
			},
		},
	},
	{
		ID:          3,
		OwnerID:     1,
		ParentTagID: cptr.Int(1),
		Description: cptr.String("test tag"),
		Names:       []string{"tag2", "tag2_alias"},
	},
	{
		ID:          8,
		OwnerID:     1,
		ParentTagID: cptr.Int(1),
		Description: cptr.String("test tag"),
		Names:       []string{"tag7", "tag7_alias"},
		Subtags: []storage.TagData{
			{
				ID:          9,
				OwnerID:     1,
				ParentTagID: cptr.Int(8),
				Description: cptr.String("test tag"),
				Names:       []string{"tag8", "tag8_alias"},
			},
		},
	},
}

func expectPath(tx *sql.Tx, si *StoragePostgres, userID int, path string, expectedID int) error {
	tagID, err := si.GetTagIDByPath(tx, userID, path)
	if err != nil {
		return errors.Annotatef(err, "getting tag id by path %q for user %d", path, userID)
	}
	if tagID != expectedID {
		return errors.Errorf(
			"GetTagIDByPath(%d, %q) should return %d, but got %d",
			userID, path, expectedID, tagID,
		)
	}
	return nil
}

func expectPathNotFound(tx *sql.Tx, si *StoragePostgres, userID int, path string) error {
	tagID, err := si.GetTagIDByPath(tx, userID, path)
	if errors.Cause(err) != storage.ErrTagDoesNotExist {
		return errors.Errorf(
			"cause of the error returned by GetTagIDByPath(%d, %q) should be ErrTagDoesNotExist (%q), but got %q, and returned id %d",
			userID, path, storage.ErrTagDoesNotExist, errors.Cause(err), tagID,
		)
	}
	return nil
}

func TestGetTagIDByPath(t *testing.T) {
	runWithRealDB(t, func(si *StoragePostgres) error {
		err := si.Tx(func(tx *sql.Tx) error {
			var u1ID, u2ID int
			var err error
			if u1ID, _, err = testutils.CreateTestUser(si, "test1", "1@1.1"); err != nil {
				return errors.Trace(err)
			}
			if u2ID, _, err = testutils.CreateTestUser(si, "test2", "2@2.2"); err != nil {
				return errors.Trace(err)
			}

			u1TagIDs, err := makeTagsHierarchy(tx, si, u1ID)
			if err != nil {
				return errors.Annotatef(err, "creating test tags hierarchy for user1")
			}

			u2TagIDs, err := makeTagsHierarchy(tx, si, u2ID)
			if err != nil {
				return errors.Annotatef(err, "creating test tags hierarchy for user2")
			}

			if err := expectPath(tx, si, u1ID, "/tag1/tag3/tag5/tag6", u1TagIDs.tag6ID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "tag1/tag3/tag5/tag6", u1TagIDs.tag6ID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "tag1/tag3_alias/tag5/tag6_alias", u1TagIDs.tag6ID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "/tag1/tag3/tag5", u1TagIDs.tag5ID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "/tag1/tag3/", u1TagIDs.tag3ID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "tag1", u1TagIDs.tag1ID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "", u1TagIDs.rootTagID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u1ID, "/", u1TagIDs.rootTagID); err != nil {
				return errors.Trace(err)
			}

			if err := expectPathNotFound(tx, si, u1ID, "/tag2/tag3"); err != nil {
				return errors.Trace(err)
			}

			if err := expectPath(tx, si, u2ID, "/tag1/tag3/tag5/tag6", u2TagIDs.tag6ID); err != nil {
				return errors.Trace(err)
			}

			return nil
		})
		return errors.Trace(err)
	})
}

func TestGetTag(t *testing.T) {
	runWithRealDB(t, func(si *StoragePostgres) error {
		var u1ID int
		var err error
		if u1ID, _, err = testutils.CreateTestUser(si, "test1", "1@1.1"); err != nil {
			return errors.Trace(err)
		}

		var rootTagID int
		var tagsData []storage.TagData

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = makeTagsHierarchy(tx, si, u1ID)
			if err != nil {
				return errors.Annotatef(err, "creating test tags hierarchy for user1")
			}

			rootTagID, err = si.GetRootTagID(tx, u1ID)
			if err != nil {
				return errors.Annotatef(err, "getting root tag for user %d", u1ID)
			}

			var err error
			tagsData, err = si.GetTags(tx, rootTagID, &storage.GetTagOpts{
				GetNames:   true,
				GetSubtags: true,
			})
			if err != nil {
				return errors.Trace(err)
			}

			if !reflect.DeepEqual(tagsData, tagsDataCreated) {
				t.Logf("%v", tagsData)
				t.Logf("%v", tagsDataCreated)
				return errors.Errorf("not equal")
			}

			return nil
		})
		if err != nil {
			return errors.Trace(err)
		}

		//panic("sdf")
		return errors.Trace(err)
	})
}

func TestInvalidTagNames(t *testing.T) {
	runWithRealDB(t, func(si *StoragePostgres) error {
		var u1ID int
		var err error
		if u1ID, _, err = testutils.CreateTestUser(si, "test1", "1@1.1"); err != nil {
			return errors.Trace(err)
		}

		var rootTagID int
		err = si.Tx(func(tx *sql.Tx) error {
			var err error
			rootTagID, err = si.GetRootTagID(tx, u1ID)
			if err != nil {
				return errors.Annotatef(err, "getting root tag for user %d", u1ID)
			}
			return nil
		})
		if err != nil {
			return errors.Trace(err)
		}

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = si.CreateTag(tx, &storage.TagData{
				OwnerID:     u1ID,
				ParentTagID: cptr.Int(rootTagID),
				Description: cptr.String("test tag"),
				Names:       []string{"123"},
			})
			return errors.Trace(err)
		})
		if err == nil || errors.Cause(err) != storage.ErrTagNameInvalid {
			return errors.Errorf("should not be able to create tag with the name 123")
		}

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = si.CreateTag(tx, &storage.TagData{
				OwnerID:     u1ID,
				ParentTagID: cptr.Int(rootTagID),
				Description: cptr.String("test tag"),
				Names:       []string{"foo bar"},
			})
			return errors.Trace(err)
		})
		if err == nil || errors.Cause(err) != storage.ErrTagNameInvalid {
			return errors.Errorf("should not be able to create tag with a space in the name")
		}

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = si.CreateTag(tx, &storage.TagData{
				OwnerID:     u1ID,
				ParentTagID: cptr.Int(rootTagID),
				Description: cptr.String("test tag"),
				Names:       []string{"foo\tbar"},
			})
			return errors.Trace(err)
		})
		if err == nil || errors.Cause(err) != storage.ErrTagNameInvalid {
			return errors.Errorf("should not be able to create tag with a tab in the name")
		}

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = si.CreateTag(tx, &storage.TagData{
				OwnerID:     u1ID,
				ParentTagID: cptr.Int(rootTagID),
				Description: cptr.String("test tag"),
				Names:       []string{"foo,bar"},
			})
			return errors.Trace(err)
		})
		if err == nil || errors.Cause(err) != storage.ErrTagNameInvalid {
			return errors.Errorf("should not be able to create tag with a comma in the name")
		}

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = si.CreateTag(tx, &storage.TagData{
				OwnerID:     u1ID,
				ParentTagID: cptr.Int(rootTagID),
				Description: cptr.String("test tag"),
				Names:       []string{string([]byte{'a', 0x01, 'b', 'c'})},
			})
			return errors.Trace(err)
		})
		if err == nil || errors.Cause(err) != storage.ErrTagNameInvalid {
			return errors.Errorf("should not be able to create tag with non-printable chars in the name")
		}

		err = si.Tx(func(tx *sql.Tx) error {
			_, err = si.CreateTag(tx, &storage.TagData{
				OwnerID:     u1ID,
				ParentTagID: cptr.Int(rootTagID),
				Names:       []string{string([]byte{'a', 0x01, 'b', 'c'})},
			})
			return errors.Trace(err)
		})
		if err == nil || errors.Cause(err) != storage.ErrTagNameInvalid {
			return errors.Errorf("should be able to create a tag without a description")
		}

		return nil
	})
}
