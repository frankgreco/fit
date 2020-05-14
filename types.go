package main

import "errors"

type File struct {
	Header  *FileHeader  `json:"header"`
	Records []DataRecord `json:"records"`
}

type FileHeader struct {
	Size            uint8  `json:"size"`
	ProtocolVersion uint8  `json:"protocol_version"`
	ProfileVersion  uint16 `json:"profile_version"`
	DataSize        uint32 `json:"data_size"`
	DataType        string `json:"data_type"`
	CRC             uint16 `json:"crc"`
}

type DataRecord struct {
	Header            *DataRecordHeader  `json:"header"`
	DefinitionMessage *DefinitionMessage `json:"definition_message,omitempty"`
	DataMessage       *DataMessage       `json:"data_message,omitempty"`
}

type DataRecordHeader struct {
	Type             DataRecordHeaderType  `json:"type"`
	LocalMessageType uint8                 `json:"local_message_type"`
	MessageType      DataRecordMessageType `json:"message_type,omitempty"`
	DeveloperData    bool                  `json:"developer_data"`
	TimeOffset       uint8                 `json:"time_offset"`
}

type FieldDefinitions []FieldDefinition

type DefinitionMessage struct {
	Architecture      uint8             `json:"architecture"`
	GlobalMessageType GlobalMessageType `json:"global_message_type"`
	NumFields         uint8             `json:"num_fields"`
	Fields            FieldDefinitions  `json:"fields"`
}

type FieldDefinition struct {
	Type               DataRecordFieldType `json:"type"`
	Number             uint8               `json:"number"`
	Size               uint8               `json:"size"`
	BaseType           *BaseType           `json:"base_type,omitempty"`
	DeveloperDataIndex uint8               `json:"developer_data_index"`
}

type BaseType struct {
	Number        uint8 `json:"number"`
	EndianAbility uint8 `json:"endian_ability"`
}

type DataMessage struct {
	NormalFields    []uint64 `json:"normal_fields"`
	DeveloperFields [][]byte `json:"developer_fields"`
}

type DataRecordHeaderType int

type DataRecordMessageType int

type DataRecordFieldType int

type GlobalMessageType int

const (
	DataRecordHeaderType_Normal DataRecordHeaderType = iota
	DataRecordHeaderType_CompressedTimestamp

	DataRecordMessageType_Definition DataRecordMessageType = iota
	DataRecordMessageType_Data

	DataRecordFieldType_Developer DataRecordFieldType = iota
	DataRecordFieldType_Normal

	GlobalMessageType_FileID GlobalMessageType = iota
	GlobalMessageType_Capabilities
	GlobalMessageType_DeviceSettings
	GlobalMessageType_UserProfile
	GlobalMessageType_HrmProfile
	GlobalMessageType_SdmProfile
	GlobalMessageType_BikeProfile
	GlobalMessageType_ZonesTarget
	GlobalMessageType_HrZone
	GlobalMessageType_PowerZone
	GlobalMessageType_MetZone
	GlobalMessageType_Sport
	GlobalMessageType_Goal
	GlobalMessageType_Session
	GlobalMessageType_Lap
	GlobalMessageType_Record
	GlobalMessageType_Event
	GlobalMessageType_DeviceInfo
	GlobalMessageType_Workout
	GlobalMessageType_WorkoutStep
	GlobalMessageType_Schedule
	GlobalMessageType_WeightScale
	GlobalMessageType_Course
	GlobalMessageType_CoursePoint
	GlobalMessageType_Totals
	GlobalMessageType_Activity
	GlobalMessageType_Software
	GlobalMessageType_FileCapabilities
	GlobalMessageType_MesgCapabilities
	GlobalMessageType_FieldCapabilites
	GlobalMessageType_FileCreator
	GlobalMessageType_BloodPressure
	GlobalMessageType_SpeedZone
	GlobalMessageType_Monitoring
	GlobalMessageType_TrainingFile
	GlobalMessageType_Hrv
	GlobalMessageType_AntRx
	GlobalMessageType_AntTx
	GlobalMessageType_AntChannelId
	GlobalMessageType_Length
	GlobalMessageType_MonitoringInfo
	GlobalMessageType_Pad
	GlobalMessageType_SlaveDevice
	GlobalMessageType_Connectivity
	GlobalMessageType_WeatherConditions
	GlobalMessageType_WeatherAlert
	GlobalMessageType_CadenceZone
	GlobalMessageType_Hr
	GlobalMessageType_SegmentLap
	GlobalMessageType_MemoGlob
	GlobalMessageType_SegmentId
	GlobalMessageType_SegmentLeaderboardEntry
	GlobalMessageType_SegmentPoint
	GlobalMessageType_SegmentFile
	GlobalMessageType_WorkoutSession
	GlobalMessageType_WatchfaceSettings
	GlobalMessageType_GpsMetadata
	GlobalMessageType_CameraEvent
	GlobalMessageType_TimestampCorrelation
	GlobalMessageType_GyroscopeData
	GlobalMessageType_AccelerometerData
	GlobalMessageType_ThreeDSensorCalibration
	GlobalMessageType_VideoFrame
	GlobalMessageType_ObdiiData
	GlobalMessageType_NmeaSentence
	GlobalMessageType_AviationAttitude
	GlobalMessageType_Video
	GlobalMessageType_VideoTitle
	GlobalMessageType_VideoDescription
	GlobalMessageType_VideoClip
	GlobalMessageType_OhrSettings
	GlobalMessageType_ExdScreenConfiguration
	GlobalMessageType_ExdDataFieldConfiguration
	GlobalMessageType_ExdDataConceptConfiguration
	GlobalMessageType_FieldDescription
	GlobalMessageType_DeveloperDataId
	GlobalMessageType_MagnetometerData
	GlobalMessageType_BarometerData
	GlobalMessageType_OneDSensorCalibration
	GlobalMessageType_Set
	GlobalMessageType_StressLevel
	GlobalMessageType_DiveSettings
	GlobalMessageType_DiveGas
	GlobalMessageType_DiveAlarm
	GlobalMessageType_ExerciseTitle
	GlobalMessageType_DiveSummary
	GlobalMessageType_Jump
	GlobalMessageType_ClimbPro
	GlobalMessageType_Unknown
	GlobalMessageType_MfgRangeMin = 0xFF00
	GlobalMessageType_MfgRangeMax = 0xFFFE
)

var (
	GlobalMessageNumber_Types = map[uint16]GlobalMessageType{
		0:   GlobalMessageType_FileID,
		1:   GlobalMessageType_Capabilities,
		2:   GlobalMessageType_DeviceSettings,
		3:   GlobalMessageType_UserProfile,
		4:   GlobalMessageType_HrmProfile,
		5:   GlobalMessageType_SdmProfile,
		6:   GlobalMessageType_BikeProfile,
		7:   GlobalMessageType_ZonesTarget,
		8:   GlobalMessageType_HrZone,
		9:   GlobalMessageType_PowerZone,
		10:  GlobalMessageType_MetZone,
		12:  GlobalMessageType_Sport,
		15:  GlobalMessageType_Goal,
		18:  GlobalMessageType_Session,
		19:  GlobalMessageType_Lap,
		20:  GlobalMessageType_Record,
		21:  GlobalMessageType_Event,
		23:  GlobalMessageType_DeviceInfo,
		26:  GlobalMessageType_Workout,
		27:  GlobalMessageType_WorkoutStep,
		28:  GlobalMessageType_Schedule,
		30:  GlobalMessageType_WeightScale,
		31:  GlobalMessageType_Course,
		32:  GlobalMessageType_CoursePoint,
		33:  GlobalMessageType_Totals,
		34:  GlobalMessageType_Activity,
		35:  GlobalMessageType_Software,
		37:  GlobalMessageType_FileCapabilities,
		38:  GlobalMessageType_MesgCapabilities,
		39:  GlobalMessageType_FieldCapabilites,
		49:  GlobalMessageType_FileCreator,
		51:  GlobalMessageType_BloodPressure,
		53:  GlobalMessageType_SpeedZone,
		55:  GlobalMessageType_Monitoring,
		72:  GlobalMessageType_TrainingFile,
		78:  GlobalMessageType_Hrv,
		80:  GlobalMessageType_AntRx,
		81:  GlobalMessageType_AntTx,
		82:  GlobalMessageType_AntChannelId,
		101: GlobalMessageType_Length,
		103: GlobalMessageType_MonitoringInfo,
		105: GlobalMessageType_Pad,
		106: GlobalMessageType_SlaveDevice,
		127: GlobalMessageType_Connectivity,
		128: GlobalMessageType_WeatherConditions,
		129: GlobalMessageType_WeatherAlert,
		131: GlobalMessageType_CadenceZone,
		132: GlobalMessageType_Hr,
		142: GlobalMessageType_SegmentLap,
		145: GlobalMessageType_MemoGlob,
		148: GlobalMessageType_SegmentId,
		149: GlobalMessageType_SegmentLeaderboardEntry,
		150: GlobalMessageType_SegmentPoint,
		151: GlobalMessageType_SegmentFile,
		158: GlobalMessageType_WorkoutSession,
		159: GlobalMessageType_WatchfaceSettings,
		160: GlobalMessageType_GpsMetadata,
		161: GlobalMessageType_CameraEvent,
		162: GlobalMessageType_TimestampCorrelation,
		164: GlobalMessageType_GyroscopeData,
		165: GlobalMessageType_AccelerometerData,
		167: GlobalMessageType_ThreeDSensorCalibration,
		169: GlobalMessageType_VideoFrame,
		174: GlobalMessageType_ObdiiData,
		177: GlobalMessageType_NmeaSentence,
		178: GlobalMessageType_AviationAttitude,
		184: GlobalMessageType_Video,
		185: GlobalMessageType_VideoTitle,
		186: GlobalMessageType_VideoDescription,
		187: GlobalMessageType_VideoClip,
		188: GlobalMessageType_OhrSettings,
		200: GlobalMessageType_ExdScreenConfiguration,
		201: GlobalMessageType_ExdDataFieldConfiguration,
		202: GlobalMessageType_ExdDataConceptConfiguration,
		206: GlobalMessageType_FieldDescription,
		207: GlobalMessageType_DeveloperDataId,
		208: GlobalMessageType_MagnetometerData,
		209: GlobalMessageType_BarometerData,
		210: GlobalMessageType_OneDSensorCalibration,
		225: GlobalMessageType_Set,
		227: GlobalMessageType_StressLevel,
		258: GlobalMessageType_DiveSettings,
		259: GlobalMessageType_DiveGas,
		262: GlobalMessageType_DiveAlarm,
		264: GlobalMessageType_ExerciseTitle,
		268: GlobalMessageType_DiveSummary,
		285: GlobalMessageType_Jump,
		317: GlobalMessageType_ClimbPro,
	}

	GlobalMessageType_Names = map[GlobalMessageType]string{
		GlobalMessageType_FileID:                      "FILE_ID",
		GlobalMessageType_Capabilities:                "CAPABILITIES",
		GlobalMessageType_DeviceSettings:              "DEVICE_SETTINGS",
		GlobalMessageType_UserProfile:                 "USER_PROFILE",
		GlobalMessageType_HrmProfile:                  "HRM_PROFILE",
		GlobalMessageType_SdmProfile:                  "SDM_PROFILE",
		GlobalMessageType_BikeProfile:                 "BIKE_PROFILE",
		GlobalMessageType_ZonesTarget:                 "ZONES_TARGET",
		GlobalMessageType_HrZone:                      "HR_ZONE",
		GlobalMessageType_PowerZone:                   "POWER_ZONE",
		GlobalMessageType_MetZone:                     "MET_ZONE",
		GlobalMessageType_Sport:                       "SPORT",
		GlobalMessageType_Goal:                        "GOAL",
		GlobalMessageType_Session:                     "SESSION",
		GlobalMessageType_Lap:                         "LAP",
		GlobalMessageType_Record:                      "RECORD",
		GlobalMessageType_Event:                       "EVENT",
		GlobalMessageType_DeviceInfo:                  "DEVICE_INFO",
		GlobalMessageType_Workout:                     "WORKOUT",
		GlobalMessageType_WorkoutStep:                 "WORKOUT_STEP",
		GlobalMessageType_Schedule:                    "SCHEDULE",
		GlobalMessageType_WeightScale:                 "WEIGHT_SCALE",
		GlobalMessageType_Course:                      "COURSE",
		GlobalMessageType_CoursePoint:                 "COURSE_POINT",
		GlobalMessageType_Totals:                      "TOTALS",
		GlobalMessageType_Activity:                    "ACTIVITY",
		GlobalMessageType_Software:                    "SOFTWARE",
		GlobalMessageType_FileCapabilities:            "FILE_CAPABILITIES",
		GlobalMessageType_MesgCapabilities:            "MESG_CAPABILITIES",
		GlobalMessageType_FieldCapabilites:            "FIELD_CAPABILITIES",
		GlobalMessageType_FileCreator:                 "FILE_CREATOR",
		GlobalMessageType_BloodPressure:               "BLOOD_PRESSURE",
		GlobalMessageType_SpeedZone:                   "SPEED_ZONE",
		GlobalMessageType_Monitoring:                  "MONITORING",
		GlobalMessageType_TrainingFile:                "TRAINING_FILE",
		GlobalMessageType_Hrv:                         "HRV",
		GlobalMessageType_AntRx:                       "ANT_RX",
		GlobalMessageType_AntTx:                       "ANT_TX",
		GlobalMessageType_AntChannelId:                "ANT_CHANNEL_ID",
		GlobalMessageType_Length:                      "LENGTH",
		GlobalMessageType_MonitoringInfo:              "MONITORING_INFO",
		GlobalMessageType_Pad:                         "PAD",
		GlobalMessageType_SlaveDevice:                 "SLAVE_DEVICE",
		GlobalMessageType_Connectivity:                "CONNECTIVITY",
		GlobalMessageType_WeatherConditions:           "WEATHER_CONDITIONS",
		GlobalMessageType_WeatherAlert:                "WEATHER_ALERT",
		GlobalMessageType_CadenceZone:                 "CADENCE_ZONE",
		GlobalMessageType_Hr:                          "HR",
		GlobalMessageType_SegmentLap:                  "SEGMENT_LAP",
		GlobalMessageType_MemoGlob:                    "MEMO_GLOB",
		GlobalMessageType_SegmentId:                   "SEGMENT_ID",
		GlobalMessageType_SegmentLeaderboardEntry:     "SEGMENT_LEADERBOARD_ENTRY",
		GlobalMessageType_SegmentPoint:                "SEGMENT_POINT",
		GlobalMessageType_SegmentFile:                 "SEGMENT_FILE",
		GlobalMessageType_WorkoutSession:              "WORKOUT_SESSION",
		GlobalMessageType_WatchfaceSettings:           "WATCHFACE_SETTINGS",
		GlobalMessageType_GpsMetadata:                 "GPS_METADATA",
		GlobalMessageType_CameraEvent:                 "CAMERA_EVENT",
		GlobalMessageType_TimestampCorrelation:        "TIMESTAMP_CORRELATION",
		GlobalMessageType_GyroscopeData:               "GYROSCOPE_DATA",
		GlobalMessageType_AccelerometerData:           "ACCELEROMETER_DATA",
		GlobalMessageType_ThreeDSensorCalibration:     "3D_SENSOR_CALIBRATION",
		GlobalMessageType_VideoFrame:                  "VIDEO_FRAME",
		GlobalMessageType_ObdiiData:                   "OBDII_DATA",
		GlobalMessageType_NmeaSentence:                "NMEA_SENTENCE",
		GlobalMessageType_AviationAttitude:            "AVIATION_ATTITUDE",
		GlobalMessageType_Video:                       "VIDEO",
		GlobalMessageType_VideoTitle:                  "VIDEO_TITLE",
		GlobalMessageType_VideoDescription:            "VIDEO_DESCRIPTION",
		GlobalMessageType_VideoClip:                   "VIDEO_CLIP",
		GlobalMessageType_OhrSettings:                 "OHR_SETTINGS",
		GlobalMessageType_ExdScreenConfiguration:      "EXD_SCREEN_CONFIGURATION",
		GlobalMessageType_ExdDataFieldConfiguration:   "EXD_DATA_FIELD_CONFIGURATION",
		GlobalMessageType_ExdDataConceptConfiguration: "EXD_DATA_CONCEPT_CONFIGURATION",
		GlobalMessageType_FieldDescription:            "FIELD_DESCRIPTION",
		GlobalMessageType_DeveloperDataId:             "DEVELOPER_DATA_ID",
		GlobalMessageType_MagnetometerData:            "MAGNETOMETER_DATA",
		GlobalMessageType_BarometerData:               "BAROMETER_DATA",
		GlobalMessageType_OneDSensorCalibration:       "1D_SENSOR_CALIBRATION",
		GlobalMessageType_Set:                         "SET",
		GlobalMessageType_StressLevel:                 "STRESS_LEVEL",
		GlobalMessageType_DiveSettings:                "DIVE_SETTINGS",
		GlobalMessageType_DiveGas:                     "DIVE_GAS",
		GlobalMessageType_DiveAlarm:                   "DIVE_ALARM",
		GlobalMessageType_ExerciseTitle:               "EXERCISE_TITLE",
		GlobalMessageType_DiveSummary:                 "DIVE_SUMMARY",
		GlobalMessageType_Jump:                        "JUMP",
		GlobalMessageType_ClimbPro:                    "CLIMB_PRO",
		GlobalMessageType_Unknown:                     "UNKNOWN",
	}

	ErrorTypeNotDefined  = errors.New("type not defined")
	ErrorMalformedBuffer = errors.New("malformed buffer")

	ContextKeyDataRecordHeader        = "DATA_RECORD_HEADER"
	ContextKeyDataRecordFieldType     = "DATA_RECORD_FIELD_TYPE"
	ContextKeyCurrentDefinitionRecord = "CURRENT_DEFINITION_RECORD"
)
