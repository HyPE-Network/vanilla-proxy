package command

import (
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// MergeAvailableCommands merges two AvailableCommands packets into one.
func MergeAvailableCommands(pk1, pk2 packet.AvailableCommands) packet.AvailableCommands {
	mergedPacket := packet.AvailableCommands{}

	// Merge and deduplicate EnumValues, returns a map for index adjustments
	EnumValues, pk1EnumValuesMap, pk2EnumValuesMap := mergeUniqueStrings(pk1.EnumValues, pk2.EnumValues)
	mergedPacket.EnumValues = EnumValues

	// Merge and deduplicate ChainedSubcommandValues
	ChainedSubcommandValues, pk1ChainedMap, pk2ChainedMap := mergeUniqueStrings(pk1.ChainedSubcommandValues, pk2.ChainedSubcommandValues)
	mergedPacket.ChainedSubcommandValues = ChainedSubcommandValues

	// Merge Enums
	Enums, pk1EnumMap, pk2EnumMap := mergeUniqueEnums(pk1.Enums, pk2.Enums, pk1EnumValuesMap, pk2EnumValuesMap)
	mergedPacket.Enums = Enums

	// Merge and deduplicate Suffixes
	Suffixes, _, _ := mergeUniqueStrings(pk1.Suffixes, pk2.Suffixes)
	mergedPacket.Suffixes = Suffixes

	// Merge ChainedSubcommands
	mergedPacket.ChainedSubcommands = mergeChainedSubcommands(pk1.ChainedSubcommands, pk2.ChainedSubcommands, pk1ChainedMap, pk2ChainedMap)

	// Merge DynamicEnums
	DynamicEnums, _, _ := mergeUniqueDynamicEnumsWithIndexMaps(pk1.DynamicEnums, pk2.DynamicEnums)
	mergedPacket.DynamicEnums = DynamicEnums

	// Merge Commands
	mergedPacket.Commands = mergeUniqueCommands(pk1.Commands, pk2.Commands, pk1ChainedMap, pk2ChainedMap, pk1EnumMap, pk2EnumMap)

	// Merge Constraints
	mergedPacket.Constraints = mergeUniqueConstraints(pk1.Constraints, pk2.Constraints, pk1EnumValuesMap, pk2EnumValuesMap)

	return mergedPacket
}

// Helper function to merge and deduplicate strings and return a map for index adjustments
func mergeUniqueStrings(slice1, slice2 []string) ([]string, map[uint]uint, map[uint]uint) {
	uniqueMap := make(map[string]uint)
	indexMap1 := make(map[uint]uint)
	indexMap2 := make(map[uint]uint)
	uniqueSlice := make([]string, 0)

	for i, item := range slice1 {
		if idx, exists := uniqueMap[item]; exists {
			// Already exists, map previous index to this.
			indexMap1[uint(i)] = idx
		} else {
			// Not found, add to uniqueSlice and map new index.
			newIdx := uint(len(uniqueSlice))
			uniqueMap[item] = newIdx
			indexMap1[uint(i)] = newIdx
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	for i, item := range slice2 {
		if idx, exists := uniqueMap[item]; exists {
			indexMap2[uint(i)] = idx
		} else {
			newIdx := uint(len(uniqueSlice))
			uniqueMap[item] = newIdx
			indexMap2[uint(i)] = newIdx
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice, indexMap1, indexMap2
}

// Helper function to merge and deduplicate CommandEnum slices, will also pass index maps for EnumValues
func mergeUniqueEnums(slice1, slice2 []protocol.CommandEnum, enumValuesIndexMap1, enumValuesIndexMap2 map[uint]uint) ([]protocol.CommandEnum, map[uint]uint, map[uint]uint) {
	uniqueMap := make(map[string]uint)
	indexMap1 := make(map[uint]uint)
	indexMap2 := make(map[uint]uint)
	uniqueSlice := make([]protocol.CommandEnum, 0)

	for i, item := range slice1 {
		item = updateEnumIndices(item, enumValuesIndexMap1)
		if idx, exists := uniqueMap[item.Type]; exists {
			// Already exists, map previous index to this.
			indexMap1[uint(i)] = idx
			// Merge the current enum, with this enum
			uniqueSlice[idx] = mergeEnums(uniqueSlice[idx], item)
		} else {
			// Not found, add to uniqueSlice and map new index.
			newIdx := uint(len(uniqueSlice))
			uniqueMap[item.Type] = newIdx
			indexMap1[uint(i)] = newIdx
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	for i, item := range slice2 {
		item = updateEnumIndices(item, enumValuesIndexMap2)
		if idx, exists := uniqueMap[item.Type]; exists {
			indexMap2[uint(i)] = idx
			// Merge the current enum, with this enum
			uniqueSlice[idx] = mergeEnums(uniqueSlice[idx], item)
		} else {
			newIdx := uint(len(uniqueSlice))
			uniqueMap[item.Type] = newIdx
			indexMap2[uint(i)] = newIdx
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice, indexMap1, indexMap2
}

// Helper function to merge two CommandEnums
func mergeEnums(enum1, enum2 protocol.CommandEnum) protocol.CommandEnum {
	uniqueIndices := make(map[uint]struct{})
	for _, idx := range enum1.ValueIndices {
		uniqueIndices[idx] = struct{}{}
	}
	for _, idx := range enum2.ValueIndices {
		uniqueIndices[idx] = struct{}{}
	}
	mergedIndices := make([]uint, 0, len(uniqueIndices))
	for idx := range uniqueIndices {
		mergedIndices = append(mergedIndices, idx)
	}
	enum1.ValueIndices = mergedIndices
	return enum1
}

// Helper function to update CommandEnum indices
func updateEnumIndices(enum protocol.CommandEnum, indexMap map[uint]uint) protocol.CommandEnum {
	updatedEnum := enum
	updatedEnum.ValueIndices = make([]uint, len(enum.ValueIndices))
	for i, idx := range enum.ValueIndices {
		updatedEnum.ValueIndices[i] = indexMap[idx]
	}
	return updatedEnum
}

// Helper function to merge ChainedSubcommand slices
func mergeChainedSubcommands(slice1, slice2 []protocol.ChainedSubcommand, indexMap1, indexMap2 map[uint]uint) []protocol.ChainedSubcommand {
	mergedSlice := make([]protocol.ChainedSubcommand, 0)
	for _, item := range slice1 {
		mergedSlice = append(mergedSlice, updateChainedSubcommandIndices(item, indexMap1))
	}
	for _, item := range slice2 {
		mergedSlice = append(mergedSlice, updateChainedSubcommandIndices(item, indexMap2))
	}
	return mergedSlice
}

// Helper function to update ChainedSubcommand indices
func updateChainedSubcommandIndices(chained protocol.ChainedSubcommand, indexMap map[uint]uint) protocol.ChainedSubcommand {
	updatedChained := chained
	updatedChained.Values = make([]protocol.ChainedSubcommandValue, len(chained.Values))
	for i, val := range chained.Values {
		val.Index = uint16(indexMap[uint(val.Index)])
		updatedChained.Values[i] = val
	}
	return updatedChained
}

// Helper function to merge and deduplicate Command slices
func mergeUniqueCommands(slice1, slice2 []protocol.Command, chainedCommandsIndexMap1, chainedCommandsIndexMap2, enumsIndexMap1, enumsIndexMap2 map[uint]uint) []protocol.Command {
	uniqueMap := make(map[string]protocol.Command)
	for _, item := range slice1 {
		uniqueMap[item.Name] = updateCommandIndices(item, chainedCommandsIndexMap1, enumsIndexMap1)
	}
	for _, item := range slice2 {
		if _, exists := uniqueMap[item.Name]; exists {
			continue // Skip duplicates
		} else {
			uniqueMap[item.Name] = updateCommandIndices(item, chainedCommandsIndexMap2, enumsIndexMap2)
		}
	}
	uniqueSlice := make([]protocol.Command, 0, len(uniqueMap))
	for _, item := range uniqueMap {
		uniqueSlice = append(uniqueSlice, item)
	}
	return uniqueSlice
}

// Helper function to update Command indices
func updateCommandIndices(cmd protocol.Command, chainedCommandsIndexMap, enumsIndexMap map[uint]uint) protocol.Command {
	updatedCmd := cmd
	updatedCmd.ChainedSubcommandOffsets = make([]uint16, len(cmd.ChainedSubcommandOffsets))
	for i, offset := range cmd.ChainedSubcommandOffsets {
		updatedCmd.ChainedSubcommandOffsets[i] = uint16(chainedCommandsIndexMap[uint(offset)])
	}
	if updatedCmd.AliasesOffset != ^uint32(0) {
		updatedCmd.AliasesOffset = uint32(enumsIndexMap[uint(cmd.AliasesOffset)])
	}
	return updatedCmd
}

// Helper function to merge and deduplicate DynamicEnum slices and return index maps
func mergeUniqueDynamicEnumsWithIndexMaps(slice1, slice2 []protocol.DynamicEnum) ([]protocol.DynamicEnum, map[uint]uint, map[uint]uint) {
	uniqueMap := make(map[string]uint)
	indexMap1 := make(map[uint]uint)
	indexMap2 := make(map[uint]uint)
	uniqueSlice := make([]protocol.DynamicEnum, 0)

	for i, item := range slice1 {
		if idx, exists := uniqueMap[item.Type]; exists {
			// Already exists, merge the values and map previous index to this.
			log.Logger.Println("Duplicate dynamic enum type in 'slice1': ", item.Type)
			indexMap1[uint(i)] = idx
			mergedValues, _, _ := mergeUniqueStrings(uniqueSlice[idx].Values, item.Values)
			uniqueSlice[idx].Values = mergedValues
		} else {
			// Not found, add to uniqueSlice and map new index.
			newIdx := uint(len(uniqueSlice))
			uniqueMap[item.Type] = newIdx
			indexMap1[uint(i)] = newIdx
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	for i, item := range slice2 {
		if idx, exists := uniqueMap[item.Type]; exists {
			log.Logger.Println("Duplicate dynamic enum type in 'slice2': ", item.Type)
			indexMap2[uint(i)] = idx
			// Merge the current dynamic enum values with this one.
			mergedValues, _, _ := mergeUniqueStrings(uniqueSlice[idx].Values, item.Values)
			uniqueSlice[idx].Values = mergedValues
		} else {
			newIdx := uint(len(uniqueSlice))
			uniqueMap[item.Type] = newIdx
			indexMap2[uint(i)] = newIdx
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice, indexMap1, indexMap2
}

// Helper function to merge and deduplicate CommandEnumConstraint slices
func mergeUniqueConstraints(slice1, slice2 []protocol.CommandEnumConstraint, indexMap1, indexMap2 map[uint]uint) []protocol.CommandEnumConstraint {
	uniqueMap := make(map[uint32]protocol.CommandEnumConstraint)
	for _, item := range slice1 {
		uniqueMap[item.EnumIndex] = updateConstraintIndices(item, indexMap1)
	}
	for _, item := range slice2 {
		if existing, exists := uniqueMap[item.EnumIndex]; exists {
			uniqueMap[item.EnumIndex] = mergeConstraints(existing, updateConstraintIndices(item, indexMap2))
		} else {
			uniqueMap[item.EnumIndex] = updateConstraintIndices(item, indexMap2)
		}
	}
	uniqueSlice := make([]protocol.CommandEnumConstraint, 0, len(uniqueMap))
	for _, item := range uniqueMap {
		uniqueSlice = append(uniqueSlice, item)
	}
	return uniqueSlice
}

// Helper function to merge two CommandEnumConstraints
func mergeConstraints(constraint1, constraint2 protocol.CommandEnumConstraint) protocol.CommandEnumConstraint {
	uniqueConstraints := make(map[byte]struct{})
	for _, c := range constraint1.Constraints {
		uniqueConstraints[c] = struct{}{}
	}
	for _, c := range constraint2.Constraints {
		uniqueConstraints[c] = struct{}{}
	}
	mergedConstraints := make([]byte, 0, len(uniqueConstraints))
	for c := range uniqueConstraints {
		mergedConstraints = append(mergedConstraints, c)
	}
	constraint1.Constraints = mergedConstraints
	return constraint1
}

// Helper function to update CommandEnumConstraint indices
func updateConstraintIndices(constraint protocol.CommandEnumConstraint, indexMap map[uint]uint) protocol.CommandEnumConstraint {
	updatedConstraint := constraint
	updatedConstraint.EnumValueIndex = uint32(indexMap[uint(constraint.EnumValueIndex)])
	updatedConstraint.EnumIndex = uint32(indexMap[uint(constraint.EnumIndex)])
	return updatedConstraint
}
